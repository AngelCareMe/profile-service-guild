package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	DB struct {
		Host    string `env:"DB_HOST" envDefault:"localhost"`
		Port    int    `env:"DB_PORT" envDefault:"5432"`
		User    string `env:"DB_USER,required"`
		Pass    string `env:"DB_PASS,required"`
		Name    string `env:"DB_NAME" envDefault:"postgres"`
		SSLMode string `env:"DB_SSL_MODE" envDefault:"disable"`
	}
	Server struct {
		Host string `env:"SERVER_HOST" envDefault:"localhost"`
		Port int    `env:"SERVER_PORT" envDefault:"8080"`
	}
	Logger struct {
		Level string `env:"LOGGER_LEVEL" envDefault:"debug"`
	}
	JWT struct {
		Secret string `env:"JWT_SECRET,notEmpty"`
	}
	Blizzard struct {
		ClientID     string `env:"BLIZZARD_CLIENT_ID,required"`
		ClientSecret string `env:"BLIZZARD_CLIENT_SECRET,required"`
		RedirectURI  string `env:"BLIZZARD_REDIRECT_URI,required"`
	}
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed parse config: %w", err)
	}

	return &cfg, nil
}
