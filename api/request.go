package api

import (
	"io"
	"net/http"
	"nft-info-collector/config"

	"github.com/kataras/golog"
	"github.com/tidwall/gjson"
)

func GetNFTScan(url string, logger *golog.Logger) string {
	// build request
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error("[API] Failed to build nftscan request")
		panic(err)
	}
	req.Header.Add("X-API-KEY", config.Load().APIKey)

	// send request
	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error("[API] Failed to send nftscan request")
		panic(err)
	}
	defer res.Body.Close()

	// analysis response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("[API] Failed to read nftscan response")
		panic(err)
	}
	return gjson.Get(string(body), "data").String()
}
