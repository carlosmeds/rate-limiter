package configs

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	DefaultLimit  int64 `mapstructure:"DEFAULT_LIMIT"`
	SecretKey     string `mapstructure:"SECRET_KEY"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
	BlockedTime   int64 `mapstructure:"BLOCKED_TIME"`
}

var (
	config *Config
	once   sync.Once
)

func LoadConfig() (*Config, error) {
	var err error
	once.Do(func() {
		viper.SetConfigName("app_config")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		viper.SetConfigFile(".env")
		viper.AutomaticEnv()
		err = viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		err = viper.Unmarshal(&config)
		if err != nil {
			panic(err)
		}
	})
	return config, err
}

func GetConfig() *Config {
	return config
}
