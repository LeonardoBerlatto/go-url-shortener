package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port               string `mapstructure:"PORT"`
	RedisURL           string `mapstructure:"REDIS_URL"`
	DynamoDBEndpoint   string `mapstructure:"DYNAMODB_ENDPOINT"`
	AWSAccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AWSRegion          string `mapstructure:"AWS_REGION"`
}

func Load() (config Config, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
