package test

import (
	"encoding/json"
	"nft-info-collector/config"
	"nft-info-collector/http"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestGetOpenSeaCollections(t *testing.T) {
	logger := iris.New().Logger()
	conf := config.Load().OpenSea
	data := http.GetOpenSeaCollections(logger, 0, conf.PageSize)

	var collections []interface{}
	err := json.Unmarshal([]byte(data), &collections)
	if err != nil {
		t.Error(err)
	}
}

func TestGetOpenSeaCollectionInfo(t *testing.T) {
	logger := iris.New().Logger()
	data := http.GetOpenSeaCollectionInfo(logger, "doodles-official")

	var collection map[string]interface{}
	err := json.Unmarshal([]byte(data), &collection)
	if err != nil {
		t.Error(err)
	}
}

// Banned by Cloudflare with 1020 - Access denied
func TestGetOpenSeaAsset(t *testing.T) {
	logger := iris.New().Logger()
	data := http.GetOpenSeaAsset(logger, "0xb47e3cd837ddf8e4c57f05d70ab865de6e193bbb", "1")

	var asset map[string]interface{}
	err := json.Unmarshal([]byte(data), &asset)
	if err != nil {
		t.Error(err)
	}
}

func TestGetOpenSeaAssets(t *testing.T) {
	logger := iris.New().Logger()
	data := http.GetOpenSeaUserAssets(logger, "0x480dd671880768D24317FA965D00f43D25868892")

	var assets []interface{}
	err := json.Unmarshal([]byte(data), &assets)
	if err != nil {
		t.Error(err)
	}
}
