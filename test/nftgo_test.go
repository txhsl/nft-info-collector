package test

import (
	"encoding/json"
	"fmt"
	"nft-info-collector/config"
	"nft-info-collector/http"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestGetNFTGoCollections(t *testing.T) {
	logger := iris.New().Logger()
	conf := config.Load().NFTGo
	data := http.GetNFTGoCollections(logger, "0", fmt.Sprint(conf.PageSize))

	var collections []interface{}
	err := json.Unmarshal([]byte(data), &collections)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNFTGoCollectionMetrics(t *testing.T) {
	logger := iris.New().Logger()
	data := http.GetNFTGoCollectionMetrics(logger, "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d")

	var metrics map[string]interface{}
	err := json.Unmarshal([]byte(data), &metrics)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNFTGoCollectionNFTs(t *testing.T) {
	logger := iris.New().Logger()
	conf := config.Load().NFTGo
	data := http.GetNFTGoCollectionNFTs(logger, "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", "0", fmt.Sprint(conf.PageSize))

	var nfts []interface{}
	err := json.Unmarshal([]byte(data), &nfts)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNFTGoNFTRarity(t *testing.T) {
	logger := iris.New().Logger()
	data := http.GetNFTGoNFTRarity(logger, "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", "4495")

	var rarity map[string]interface{}
	err := json.Unmarshal([]byte(data), &rarity)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNFTGoNFTMetrics(t *testing.T) {
	logger := iris.New().Logger()
	data := http.GetNFTGoNFTMetrics(logger, "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", "4495")

	var metrics map[string]interface{}
	err := json.Unmarshal([]byte(data), &metrics)
	if err != nil {
		t.Error(err)
	}
}
