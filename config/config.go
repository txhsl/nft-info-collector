package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	NFTScan map[string]string `yaml:"nftscan"`
	OpenSea map[string]int    `yaml:"opensea"`
	MongoDB map[string]string `yaml:"mongodb"`
	Port    string            `yaml:"port"`
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
