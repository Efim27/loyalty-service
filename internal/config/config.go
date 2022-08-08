package config

import (
	"flag"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddr  string `mapstructure:"RUN_ADDRESS"`
	DBSource    string `mapstructure:"DATABASE_URI"`
	AccrualAddr string `mapstructure:"ACCRUAL_SYSTEM_ADDRESS"`
}

func (config *Config) LoadEnv(path string) (err error) {
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

func (config *Config) LoadFlags() {
	flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "server address (host:port)")
	flag.StringVar(&config.DBSource, "d", config.DBSource, "Postgres DSN")
	flag.StringVar(&config.AccrualAddr, "r", config.AccrualAddr, "Accrual addr")

	flag.Parse()
}

func LoadConfig() (config Config, err error) {
	config.LoadFlags()
	err = config.LoadEnv("./config")

	return
}
