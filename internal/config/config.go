package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddr  string `mapstructure:"RUN_ADDRESS"`
	DBSource    string `mapstructure:"DATABASE_URI"`
	AccrualAddr string `mapstructure:"ACCRUAL_SYSTEM_ADDRESS"`
}

func LoadEnv(config *Config, path string) (err error) {
	viper.AllowEmptyEnv(true)

	viper.SetDefault("RUN_ADDRESS", "127.0.0.1:8081")
	viper.SetDefault("ACCRUAL_SYSTEM_ADDRESS", "127.0.0.1:8080")

	viper.AddConfigPath(path)
	viper.SetConfigName("main")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func LoadConfig(path string) (config Config, err error) {
	LoadEnv(&config, "./config")

	return
}
