package configs

import (
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	DefaultLimit  int64  `mapstructure:"DEFAULT_LIMIT"`
	SecretKey     string `mapstructure:"SECRET_KEY"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
	BlockedTime   int64  `mapstructure:"BLOCKED_TIME"`
	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	ApiKeyLimits  map[string]int64
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

		config.ApiKeyLimits = make(map[string]int64)
		apiKeys := viper.GetString("API_KEYS")
		for _, pair := range strings.Split(apiKeys, ",") {
			parts := strings.Split(pair, ":")
			if len(parts) == 2 {
				apiKey := parts[0]
				limit, err := strconv.ParseInt(parts[1], 10, 64)
				if err == nil {
					config.ApiKeyLimits[apiKey] = limit
				}
			}
		}
	})
	return config, err
}

func GetConfig() *Config {
	return config
}
