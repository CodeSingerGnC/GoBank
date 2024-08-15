package main

import (
	// base
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	// self
	"github.com/CodeSingerGnC/MicroBank/api"
	db "github.com/CodeSingerGnC/MicroBank/db/sqlc"
	_ "github.com/CodeSingerGnC/MicroBank/doc/statik"
	"github.com/CodeSingerGnC/MicroBank/gapi"
	"github.com/CodeSingerGnC/MicroBank/mail"
	"github.com/CodeSingerGnC/MicroBank/pb"
	"github.com/CodeSingerGnC/MicroBank/util"
	"github.com/CodeSingerGnC/MicroBank/worker"

	// static
	"github.com/rakyll/statik/fs"

	// grpc
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql dirver
	_ "github.com/go-sql-driver/mysql"
	// -TODO: golong-migrate
	// "github.com/golang-migrate/migrate/v4"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	// -TODO: run db migration

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisServerAddress,
	}

	waitGroup, ctx := errgroup.WithContext(ctx)

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	// 如果需要使用 GinServer，请注释 runGatewayServer 并取消 runGinServer 的注释。
	runTaskProcessor(ctx, waitGroup, redisOpt, store, config)
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor)
	// runGinServer(config, store)
	runGprcServer(ctx, waitGroup, config, store, taskDistributor)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runTaskProcessor(
	ctx context.Context,
	waitGroup *errgroup.Group,
	redisOpt asynq.RedisClientOpt, 
	store db.Store, 
	config util.Config,
) {
	mailer := mail.NewSinaSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)

	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

	waitGroup.Go(func () error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown task processor")

		taskProcessor.Shutdown()
		log.Info().Msg("task processor is stopped")
		return nil
	})
}

func runGprcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,  
	config util.Config, 
	store db.Store, 
	taskDistributor worker.TaskDistributor,
){
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterMicroBankServer(grpcServer, server)
	// -TODO: reflection
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	waitGroup.Go(func () error {
		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("cannot start gRPC server")
			return err
		}	

		return nil
	})

	waitGroup.Go(func () error {
		<-ctx.Done()
		log.Info().Msg("graceful shout down grpc server")
		grpcServer.GracefulStop()
		log.Info().Msg("grpc server is stopped")

		return nil
	})
}

func runGatewayServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config, 
	store db.Store,
	taskDistributor worker.TaskDistributor,
){
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	grpcMux := runtime.NewServeMux()
	
	err = pb.RegisterMicroBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}
	
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik file system")
	}
	
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	httpServer := &http.Server{
		Handler: gapi.HttpLogger(mux),
		Addr: config.HTTPServerAdress,
	}

	waitGroup.Go(func () error {
		log.Info().Msgf("start HTTP gateway server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("Http gateway Server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func () error {
		<-ctx.Done()
		log.Info().Msg("graceful shout down http gateway server")
		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown http gatewat server")
			return err
		}
		log.Info().Msg("http gateway server is stopped")
		return nil
	})
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAdress)
	if err != nil {
		log.Fatal().Msg("cannot start http server")
	}
}
