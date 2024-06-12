package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	Address          string `env:"RUN_ADDRESS"`
	RedisAddress     string `env:"REDIS_URI"`
	DefaultDBConnStr string `env:"DATABASE_URI"`
}

const (
	defaultServer       = "localhost:8080"
	defaultRedisAddress = "localhost:6379"
	defaultConnStr      = "host=localhost port=5432 user=Mr8 password=Rammstein12! dbname=rutube_db sslmode=disable"
)

func New() (*Config, error) {
	cfg := Config{}

	flag.StringVar(&cfg.Address, "a", defaultServer, "server address [host:port]")
	flag.StringVar(&cfg.RedisAddress, "m", defaultRedisAddress, "redis connection string")
	flag.StringVar(&cfg.DefaultDBConnStr, "d", defaultConnStr, "connection string")

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
