package config

import (
	"github.com/kelseyhightower/envconfig"
)

func GetConfig() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return cfg, err
}

type Config struct {
	Address    string `envconfig:"GRPC_ADDRESS" default:"127.0.0.1:50051"`
	WindowSize int    `envconfig:"WINDOW_SIZE" default:"512"`
}
