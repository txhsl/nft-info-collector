package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type nftscan struct {
	ApiKey string `yaml:"api-key"`
}

type nftgo struct {
	ApiKey   string `yaml:"api-key"`
	PageSize int    `yaml:"page-size"`
	Limit    int    `yaml:"limit"`
}

type reservoir struct {
	ApiKey   string `yaml:"api-key"`
	PageSize int    `yaml:"page-size"`
	Limit    int    `yaml:"limit"`
}

type opensea struct {
	ApiKey   string `yaml:"api-key"`
	PageSize int    `yaml:"page-size"`
	Limit    int    `yaml:"limit"`
}

type mongodb struct {
	Url string `yaml:"url"`
}

type keywords struct {
	Sorts []string `yaml:"sorts"`
	Times []string `yaml:"times"`
}

type Config struct {
	NFTScan   nftscan   `yaml:"nftscan"`
	NFTGo     nftgo     `yaml:"nftgo"`
	Reservoir reservoir `yaml:"reservoir"`
	OpenSea   opensea   `yaml:"opensea"`
	MongoDB   mongodb   `yaml:"mongodb"`
	Port      string    `yaml:"port"`
	Keywords  keywords  `yaml:"keywords"`
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
