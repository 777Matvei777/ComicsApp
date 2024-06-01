package config

import (
	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Url              string `yaml:"source_url"`
	DbFile           string `yaml:"DbFile"`
	Parallel         int    `yaml:"parallel"`
	Port             string `yaml:"port"`
	Postgresql       string `yaml:"postgresql"`
	Token_max_time   int    `yaml:"token_max_time"`
	ConcurrencyLimit int    `yaml:"concurrencyLimit"`
	RateLimit        int    `yaml:"rateLimit"`
}

func New(config string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(config, &cfg); err != nil {
		return nil, errors.New("can't read config: " + err.Error())
	}
	return &cfg, nil
}
