package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config 用于储存应用程序的所有的配置，
// 其值由 viper 从配置文件或环境变量中读取。
type Config struct {
	DBDriver             	string        `mapstructure:"DB_DRIVER"`
	DBSource             	string        `mapstructure:"DB_SOURCE"`
	// MigrationURL		 	string        `mapstructure:"MIGRATIONS_URL"`
	Environment			 	string			`mapstructure:"ENVIRONMENT"`
	HTTPServerAdress     	string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    	string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	RedisServerAddress   	string        `mapstructure:"REDIS_ADDRESS"`
	TokenSymmetricKey    	string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	EmailSenderName	    	string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress    	string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword    	string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	AccessTokenDuration  	time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration 	time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig 从配置文件或应用程序中读取环境变量。
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
