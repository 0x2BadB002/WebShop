package config

import (
	"time"

	"github.com/spf13/viper"
)

type MongoCfg struct {
	Uri     string        `mapstructure:"uri"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type AuthCfg struct {
	PrivateKey string `mapstructure:"private_key"`
}

type Config struct {
	Mongo             MongoCfg      `mapstructure:"mongo"`
	Auth              AuthCfg       `mapstructure:"auth"`
	RecentItemsPeriod time.Duration `mapstructure:"recent_items_period"`
	RecentUsersCount  int64         `mapstructure:"recent_users_count"`
}

func Get() (*Config, error) {
	config := new(Config)

	viper.SetConfigType("hcl")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("mongo.uri", "mongodb://localhost:27017")
	viper.SetDefault("mongo.timeout", "10s")

	viper.SetDefault("auth.esdca", "")

	viper.SetDefault("recent_items_period", "72h")
	viper.SetDefault("recent_users_count", 2)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
