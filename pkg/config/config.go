package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Url      string `yaml:"source_url"`
	Db_file  string `yaml:"db_file"`
	Parallel int    `yaml:"parallel"`
}

func New(config string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(config, &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}
	return &cfg
}
