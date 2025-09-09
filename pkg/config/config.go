package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`
	DB struct {
		User    string `mapstructure:"user"`
		Pass    string `mapstructure:"pass"`
		Host    string `mapstructure:"host"`
		Port    int    `mapstructure:"port"`
		Name    string `mapstructure:"name"`
		SSLMode string `mapstructure:"sslmode"`
	} `mapstructure:"db"`
	Logger struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logger"`
	Blizzard struct {
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		RedirectURL  string `mapstructure:"redirect_url"`
	} `mapstructure:"blizzard"`
	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.AddConfigPath(".")
	v.AddConfigPath("..")

	v.SetConfigType("yaml")
	v.SetConfigName("config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed reading config: %w", err)
	}

	v.AutomaticEnv()

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed mapping config to struct: %w", err)
	}

	log.Printf("Loaded config: DB_HOST=%s, DB_PORT=%d, DB_USER=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User)

	return &cfg, nil
}
