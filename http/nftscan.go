package http

import (
	"io"
	"net/http"
	"nft-info-collector/config"

	"github.com/kataras/golog"
	"github.com/tidwall/gjson"
)

// Used in DB cache and immediate response
func GetNFTScanTrends(logger *golog.Logger) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://restapi.nftscan.com/api/v2/statistics/ranking/trade?time=7d&sort_field=volume&sort_direction=desc&show_7d_trends=false"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().NFTScan.ApiKey)

	// send request
	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// analysis response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return gjson.Get(string(body), "data").String(), nil
}
