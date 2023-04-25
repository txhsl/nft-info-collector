package test

import (
	"encoding/json"
	"nft-info-collector/http"
	"testing"

	"github.com/kataras/iris/v12"
)

func TestGetReservoirCollections(t *testing.T) {
	logger := iris.New().Logger()
	contracts := []string{"0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", "0x8d04a8c79ceb0889bdd12acdf3fa9d207ed3ff63", "0xb47e3cd837ddf8e4c57f05d70ab865de6e193bbb"}
	data, err := http.GetReservoirCollections(logger, contracts)
	if err != nil {
		t.Error(err)
	}

	var collections []interface{}
	err = json.Unmarshal([]byte(data), &collections)
	if err != nil {
		t.Error(err)
	}
}
