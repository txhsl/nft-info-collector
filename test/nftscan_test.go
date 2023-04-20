package test

import (
	"encoding/json"
	"nft-info-collector/http"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestGetNFTScanTrends(t *testing.T) {
	logger := iris.New().Logger()
	data := http.GetNFTScanTrends(logger)

	var collections []interface{}
	err := json.Unmarshal([]byte(data), &collections)
	if err != nil {
		t.Error(err)
	}
}
