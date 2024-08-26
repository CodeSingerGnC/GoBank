# 使用异步任务队列 asynq 发送邮件

## asynq client

asynq client 负责创建任务，设置任务选项，然后将任务放入队列。

在项目中，上述任务由封装的 TaskDistributor 接口抽象，RedisTaskDistributor 结构体实现了 TaskDistributor，所以实际任务在 RedisTaskDistributor 中定义，NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor 函数负责创建 RedisTaskDistributor 然后返回 TaskDistributor。

```go
type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PalyloadSendVerifyEmail,
		opt ...asynq.Option,
	) error 
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
```

TaskDistributor 中的 DistributeTaskSendVerifyEmail 则是实际负责将创建任务，设置任务，任务入队的函数。

其中的 PalyloadSendVerifyEmail 结构如下：

```go
type PalyloadSendVerifyEmail struct {
	UserAccount string `json:"user_account"`
}
```

只存放用户账号，在处理信息的时候，则使用用户账号查找用户邮箱，然后再根据邮箱发送邮件。

RedisTaskDistributor 对于 DistributeTaskSendVerifyEmail 的实现如下：

```go
const TaskSendVerifyEmail = "task:send_verify_email"

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PalyloadSendVerifyEmail,
	opt ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opt...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("faided to enqueue task: %w", err)
	}

	log.Info().Str("type: ", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queque", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueue task")
	return nil
}
```

此处的 TaskSendVerifyEmail 用于标记任务类型，在服务端可以通过该任务类型调用具体的函数进行处理。

这里的主要逻辑如下：

1. 将 payload 序列化为 json 格式。
2. 使用 async.NewTask 创建任务，创建时指定任务类型，任务选项。
3. 将任务入队，入队后会返回 info。info 会返回任务的基本信息。
4. 创建任务日志。

为了在方便调用 DistributeTaskSendVerifyEmail，我将 TaskDistributor 内嵌到 server(银行服务器的服务器对象) 中，在创建服务器时，会使用 NewRedisTaskProcessor 设置 server 中的 TaskDistributor(主要是根据指定的 redis 地址，创建 redis 客户端)。然后 DistributeTaskSendVerifyEmail 的调用是在创建用户时发生的，即调用 CreateUser(ctx context.Context, req *pb.CreateUserRequest) 时，因为创建用户的逻辑是在数据库创建用户并发送验证邮件，为了确保两者的原子性，使用事务操作，所以在 CreateUser 中调用 CreateUserTx，CreateUserTx 实现如下：

```go
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParam) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		err = arg.AfterCreate()
		return err
	})
	if err != nil {
		return err
	}

	return nil
}
```

而调用时的 arg 如下：

```go
	arg := db.CreateUserTxParam{
		CreateUserParams: db.CreateUserParams{
			UserAccount: req.GetUserAccount(),
			HashPassword: hashedPassword,
			Username: req.GetUsername(),
			Email: req.GetEmail(),
		},
		AfterCreate: func () error{
			// -TODO: use db transaction
			taskPayload := &worker.PalyloadSendVerifyEmail{
				UserAccount: req.GetEmail(),
			}
			opt := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opt...)
		},
	}
```

这里的 AfterCreate 为回调函数，在 CreateUserTx 中创建完用户之后会执行。

还需要注意的是这里 opt 的设置，极为 asynq 中任务的设置，这里设置了任务的最大重试次数，相对于当前时间处理给定任务的时间，加入的队列名称。

worker 中存在两个队列：

```go
const (
	QueueCritical = "critical"
	QueueDefault = "default"
)
```

这里指定的是第一个队列。

### asynq server

asynq server 负责从任务队列取出任务，然后根据任务类型调用处理函数处理任务。

async server 的操作主要包括运行服务，处理任务，关闭服务，这些操作被封装在 TaskProcessor 接口中：

```go
type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}
```

本项目中主要使用 asynq 的目的是异步发送邮件，提高响应速度（也就是不用等待邮件发送再返回），并使用 goroutine 并发处理发送邮件的任务。在操作的过程中，会涉及到邮件的发送，数据库操作，当然还需要 async server，所以我封装了 RedisTaskProcessor：

```go
type RedisTaskProcessor struct {
	server 	*asynq.Server
	store 	db.Store
	mailer 	mail.EmailSender
}
```

RedisTaskProcessor 实现了 TaskProcessor 接口的所有方法。

start 的实现：

```go
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
```

主要逻辑：

1. 创建多路复用器
2. 绑定 TaskSendVerifyEmail 任务类型和 ProcessTaskSendVerifyEmail 方法
3. 启动服务

Shutdown() 包装了 asynq.server.Shutdown()。

ProcessTaskSendVerifyEmail 实现如下：

```go
func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PalyloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.UserAccount)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	_, err = processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		UserAccount: user.UserAccount,
		Email: user.Email,
		SecretCode: util.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}
	
	log.Info().
	Str("type: ", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")
	
	return nil
}
```

主要逻辑：

1. 将任务的 payload 反序列化为 PalyloadSendVerifyEmail
2. 通过 PalyloadSendVerifyEmail 中的 userAccount 访问数据库获取信息
3. 发送邮件。
4. 在数据库创建验证邮件对象。
5. 记录日志。

那么在使用 RedisTaskProcessor 时，首先得创建 RedisTaskProcessor，其方式便是调用 NewRedisTaskProcessor：

```
func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault: 5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Info().Err(err).
					Str("type", task.Type()).
					Msg("process task failed")
			}),
			Logger: NewLogger(),
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
	}
}
```

NewRedisTaskProcessor 主要是创建 asynq server 并包装 store 和 mailer 为 RedisTaskProcessor 对象。这里创建 server 的时候存在两个字段，redisOpt，asynq.Config。redisOpt 可以指定连接方式，包括 tcp 和 unix，以及使用 TLS，最重要的是指定 redis 地址。asynq.Config 主要是配置 server，这里存在三个配置 Queues，ErrorHandler，Logger。Queues 用于指定队列处理的优先级，ErrorHandler 指定任务处理出错时的处理方式，这里主要是记录日志，没有额外的操作。Logger 用于指定日志系统，我们使用的是第三方日志  "github.com/rs/zerolog"，但是因为需要实现 asynq.Logger 的接口才可以被当成 asynq.Logger 所以我们简单地进行了封装。

Queues 的优先级具体描述如下：

```go
其中 async.Config 的 Queues 是个 map 代表的是需要处理的队列与优先级

示例：

Queues: map[string]int{

​        "critical": 6,

​         "default":  3,

​         "low":      1,

}

根据上述配置，并假设所有队列都不为空，"critical"、"default"、"low" 队列中的任务

分别应被处理 60%、30%、10% 的时间。
```

在项目中，我们包装了 runTaskProcessor 来启动服务：

```go
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
```

那么最后一步，就是在 main.go 中调用该函数启动服务了，详细可以看 main.go 文件。
