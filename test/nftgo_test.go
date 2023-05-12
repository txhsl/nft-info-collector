package test

import (
	"encoding/json"
	"nft-info-collector/config"
	"nft-info-collector/http"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/tidwall/gjson"
)

func TestGetNFTGoCollections(t *testing.T) {
	logger := iris.New().Logger()
	conf := config.Load().NFTGo
	data, err := http.GetNFTGoCollections(logger, "7d", 0, conf.PageSize)
	if err != nil {
		t.Error(err)
	}

	var collections []interface{}
	err = json.Unmarshal([]byte(gjson.Get(data, "collections").String()), &collections)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNFTGoCollectionInfo(t *testing.T) {
	logger := iris.New().Logger()
	data, err := http.GetNFTGoCollectionInfo(logger, "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d")
	if err != nil {
		t.Error(err)
	}

	var collection map[string]interface{}
	err = json.Unmarshal([]byte(data), &collection)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNFTGoCollectionMetrics(t *testing.T) {
	logger := iris.New().Logger()
	data, err := http.GetNFTGoCollectionMetrics(logger, "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d")
	if err != nil {
		t.Error(err)
	}

	var metrics map[string]interface{}
	err = json.Unmarshal([]byte(data), &metrics)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNFTGoCollectionNFTs(t *testing.T) {
	logger := iris.New().Logger()
	conf := config.Load().NFTGo
	data, err := http.GetNFTGoCollectionNFTs(logger, "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", 0, conf.PageSize)
	if err != nil {
		t.Error(err)
	}

	var nfts []interface{}
	err = json.Unmarshal([]byte(gjson.Get(data, "nfts").String()), &nfts)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNFTGoNFTRarity(t *testing.T) {
	logger := iris.New().Logger()
	data, err := http.GetNFTGoNFTRarity(logger, "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", 4495)
	if err != nil {
		t.Error(err)
	}

	var rarity map[string]interface{}
	err = json.Unmarshal([]byte(data), &rarity)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNFTGoNFTMetrics(t *testing.T) {
	logger := iris.New().Logger()
	data, err := http.GetNFTGoNFTMetrics(logger, "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", 4495)
	if err != nil {
		t.Error(err)
	}

	var metrics map[string]interface{}
	err = json.Unmarshal([]byte(data), &metrics)
	if err != nil {
		t.Error(err)
	}
}

func TestGetNFTGoUserAssets(t *testing.T) {
	logger := iris.New().Logger()
	data, err := http.GetNFTGoUserAssets(logger, "0x480dd671880768D24317FA965D00f43D25868892", 0, 10)
	if err != nil {
		t.Error(err)
	}

	var assets []interface{}
	err = json.Unmarshal([]byte(gjson.Get(data, "assets").String()), &assets)
	if err != nil {
		t.Error(err)
	}
}
