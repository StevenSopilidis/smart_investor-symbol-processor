package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Environment       string `mapstructure:"ENVIRONMENT""`
	ServiceName       string `mapstructure:"SERVICE_NAME"`
	OtelCollectorAddr string `mapstructure:"OTEL_COLLECTOR_ADDR"`
	SymbolTopic       string `mapstructure:"SYMBOL_TOPIC"`
	KafkaBroker       string `mapstructure:"KAFKA_BROKER"`
	KafkaGroupId      string `mapstructure:"KAFKA_GROUP_ID"`
	RedisAddr         string `mapstructure:"REDIS_ADDR"`
	RedisPass         string `mapstructure:"REDIS_PASSWORD"`
	RedisDB           int    `mapstructure:"REDIS_DB"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
	}

	viper.AutomaticEnv()
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("OTEL_COLLECTOR_ADDR")
	viper.BindEnv("SERVICE_NAME")
	viper.BindEnv("SYMBOL_TOPIC")
	viper.BindEnv("KAFKA_BROKER")
	viper.BindEnv("REDIS_ADDR")
	viper.BindEnv("REDIS_PASSWORD")
	viper.BindEnv("REDIS_DB")
	viper.BindEnv("KAFKA_GROUP_ID")

	err = viper.Unmarshal(&config)
	return config, err
}
