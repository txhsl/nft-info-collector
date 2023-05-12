package test

import (
	"encoding/json"
	"nft-info-collector/http"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/tidwall/gjson"
)

func TestGetNFTScanTrends(t *testing.T) {
	logger := iris.New().Logger()
	data, err := http.GetNFTScanTrends(logger)
	if err != nil {
		t.Error(err)
	}

	var collections []interface{}
	err = json.Unmarshal([]byte(gjson.Get(data, "data").String()), &collections)
	if err != nil {
		t.Error(err)
	}
}
