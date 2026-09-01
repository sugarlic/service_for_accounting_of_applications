package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Load() *Config {
	_ = godotenv.Load() // optional

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("config: parse env: %v", err)
	}

	return &cfg
}
