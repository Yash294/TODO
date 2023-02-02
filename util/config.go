package util

import (
	"github.com/spf13/viper"
)

type DBConfig struct {
	DBHost string `mapstructure:"DB_HOST"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASS"`
	DBName string `mapstructure:"DB_NAME"`
	DBPort int `mapstructure:"DB_PORT"`
}

type RedisConfig struct {
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisUser string `mapstructure:"REDIS_USER"`
	RedisPass string `mapstructure:"REDIS_PASS"`
	RedisName string `mapstructure:"REDIS_NAME"`
	RedisPort int `mapstructure:"REDIS_PORT"`
}

func LoadDBConfig(path string) (config DBConfig, err error) {
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

func LoadRedisConfig(path string) (config RedisConfig, err error) {
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