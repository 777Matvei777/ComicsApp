package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Url     string `yaml:"source_url"`
	Db_file string `yaml:"db_file"`
}

func New() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("config.yaml", &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}
	return &cfg
}
