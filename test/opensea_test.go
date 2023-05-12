package test

import (
	"encoding/json"
	"nft-info-collector/config"
	"nft-info-collector/http"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/tidwall/gjson"
)

func TestGetOpenSeaCollections(t *testing.T) {
	logger := iris.New().Logger()
	conf := config.Load().OpenSea
	data, err := http.GetOpenSeaCollections(logger, 0, conf.PageSize)
	if err != nil {
		t.Error(err)
	}

	var collections []interface{}
	err = json.Unmarshal([]byte(gjson.Get(data, "collections").String()), &collections)
	if err != nil {
		t.Error(err)
	}
}

func TestGetOpenSeaCollectionInfo(t *testing.T) {
	logger := iris.New().Logger()
	data, err := http.GetOpenSeaCollectionInfo(logger, "doodles-official")
	if err != nil {
		t.Error(err)
	}

	var collection map[string]interface{}
	err = json.Unmarshal([]byte(gjson.Get(data, "collection").String()), &collection)
	if err != nil {
		t.Error(err)
	}
}

func TestGetOpenSeaCollectionOffers(t *testing.T) {
	logger := iris.New().Logger()
	data, err := http.GetOpenSeaCollectionOffers(logger, "doodles-official")
	if err != nil {
		t.Error(err)
	}

	var offers []interface{}
	err = json.Unmarshal([]byte(gjson.Get(data, "offers").String()), &offers)
	if err != nil {
		t.Error(err)
	}
}

func TestGetOpenSeaCollectionRecentSales(t *testing.T) {
	logger := iris.New().Logger()
	data, err := http.GetOpenSeaCollectionRecentSales(logger, "doodles-official")
	if err != nil {
		t.Error(err)
	}

	var sales []interface{}
	err = json.Unmarshal([]byte(gjson.Get(data, "asset_events").String()), &sales)
	if err != nil {
		t.Error(err)
	}
}

// Banned by Cloudflare with 1020 - Access denied
func TestGetOpenSeaAsset(t *testing.T) {
	logger := iris.New().Logger()
	data, err := http.GetOpenSeaAsset(logger, "0xb47e3cd837ddf8e4c57f05d70ab865de6e193bbb", 1)
	if err != nil {
		t.Error(err)
	}

	var asset map[string]interface{}
	err = json.Unmarshal([]byte(data), &asset)
	if err != nil {
		t.Error(err)
	}
}

func TestGetOpenSeaAssets(t *testing.T) {
	logger := iris.New().Logger()
	data, err := http.GetOpenSeaUserAssets(logger, "0x480dd671880768D24317FA965D00f43D25868892", "")
	if err != nil {
		t.Error(err)
	}

	var assets []interface{}
	err = json.Unmarshal([]byte(gjson.Get(data, "assets").String()), &assets)
	if err != nil {
		t.Error(err)
	}
}
