package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type nftscan struct {
	ApiKey string `yaml:"api-key"`
}

type nftgo struct {
	ApiKey     string   `yaml:"api-key"`
	PageSize   int      `yaml:"page-size"`
	Limit      int      `yaml:"limit"`
	Keywords   []string `yaml:"keywords"`
	Filters    []string `yaml:"filters"`
	TimeRanges []string `yaml:"time-ranges"`
}

type reservoir struct {
	ApiKey   string `yaml:"api-key"`
	PageSize int    `yaml:"page-size"`
	Limit    int    `yaml:"limit"`
}

type opensea struct {
	PageSize int `yaml:"page-size"`
	Limit    int `yaml:"limit"`
}

type mongodb struct {
	Url string `yaml:"url"`
}

type Config struct {
	NFTScan   nftscan   `yaml:"nftscan"`
	NFTGo     nftgo     `yaml:"nftgo"`
	Reservoir reservoir `yaml:"reservoir"`
	OpenSea   opensea   `yaml:"opensea"`
	MongoDB   mongodb   `yaml:"mongodb"`
	Port      string    `yaml:"port"`
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
