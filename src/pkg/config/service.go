package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

func NewConfig() (*Config, error) {
	config := &Config{}

	if err := cleanenv.ReadConfig("configs/config.yml", config); err != nil {
		return nil, err
	}

	return config, nil
}
