package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Host               string `mapstructure:"HOST"`
	Port               string `mapstructure:"PORT"`
	RedisURL           string `mapstructure:"REDIS_URL"`
	DynamoDBEndpoint   string `mapstructure:"DYNAMODB_ENDPOINT"`
	AWSAccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AWSRegion          string `mapstructure:"AWS_REGION"`
	LogLevel           string `mapstructure:"LOG_LEVEL"`
}

func Load() (config Config, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("LOG_LEVEL", "info")

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
