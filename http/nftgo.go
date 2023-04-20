package http

import (
	"fmt"
	"io"
	"net/http"
	"nft-info-collector/config"

	"github.com/kataras/golog"
	"github.com/tidwall/gjson"
)

func GetNFTGoCollections(logger *golog.Logger, offset int, limit int) string {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/market/rank/collection/7d?by=volume&with_rarity=false&asc=false&offset=" + fmt.Sprint(offset) + "&limit=" + fmt.Sprint(limit)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error("[API] Failed to build nftgo request")
		panic(err)
	}
	req.Header.Add("X-API-KEY", config.Load().NFTGo.ApiKey)

	// send request
	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error("[API] Failed to send nftgo request")
		panic(err)
	}
	defer res.Body.Close()

	// analysis response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("[API] Failed to read nftgo response")
		panic(err)
	}
	return gjson.Get(string(body), "collections").String()
}
