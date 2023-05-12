package test

import (
	"encoding/json"
	"nft-info-collector/http"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/tidwall/gjson"
)

func TestGetReservoirCollections(t *testing.T) {
	logger := iris.New().Logger()
	contracts := []string{"0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", "0x8d04a8c79ceb0889bdd12acdf3fa9d207ed3ff63", "0xb47e3cd837ddf8e4c57f05d70ab865de6e193bbb"}
	data, err := http.GetReservoirCollections(logger, contracts)
	if err != nil {
		t.Error(err)
	}

	var collections []interface{}
	err = json.Unmarshal([]byte(gjson.Get(data, "collections").String()), &collections)
	if err != nil {
		t.Error(err)
	}
}

func TestGetReservoirCollectionListing(t *testing.T) {
	logger := iris.New().Logger()
	contract := "0x8d04a8c79ceb0889bdd12acdf3fa9d207ed3ff63"
	data, err := http.GetReservoirCollectionListing(logger, contract)
	if err != nil {
		t.Error(err)
	}

	var offers []interface{}
	err = json.Unmarshal([]byte(gjson.Get(data, "orders").String()), &offers)
	if err != nil {
		t.Error(err)
	}
}

func TestGetReservoirCollectionDaily(t *testing.T) {
	logger := iris.New().Logger()
	contract := "0x8d04a8c79ceb0889bdd12acdf3fa9d207ed3ff63"
	data, err := http.GetReservoirCollectionDaily(logger, contract)
	if err != nil {
		t.Error(err)
	}

	var graph []interface{}
	err = json.Unmarshal([]byte(gjson.Get(data, "collections").String()), &graph)
	if err != nil {
		t.Error(err)
	}
}
