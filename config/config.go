package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey  string `yaml:"apikey"`
	MongoDB string `yaml:"mongodb"`
	Port    string `yaml:"port"`
}

func Load() Config {
	config := Config{}
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	return config
}
