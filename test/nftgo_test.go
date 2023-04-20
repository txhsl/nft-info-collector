package test

import (
	"encoding/json"
	"nft-info-collector/config"
	"nft-info-collector/http"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestGetNFTGoCollections(t *testing.T) {
	logger := iris.New().Logger()
	conf := config.Load().NFTGo
	data := http.GetNFTGoCollections(logger, 0, conf.PageSize)

	var collections []interface{}
	err := json.Unmarshal([]byte(data), &collections)
	if err != nil {
		t.Error(err)
	}
}
