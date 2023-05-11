package http

import (
	"io"
	"net/http"
	"nft-info-collector/config"

	"github.com/kataras/golog"
	"github.com/tidwall/gjson"
)

// Only used in DB cache
// @deprecated
func GetReservoirCollections(logger *golog.Logger, contracts []string) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://api.reservoir.tools/collections/v5?includeTopBid=true"
	for _, contract := range contracts {
		url += "&contract=" + contract
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().Reservoir.ApiKey)

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

	return gjson.Get(string(body), "collections").String(), nil
}

func GetReservoirCollectionListing(logger *golog.Logger, contract string) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://api.reservoir.tools/orders/asks/v4?contracts=" + contract + "&source=opensea.io&&sortBy=price&limit=5"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().Reservoir.ApiKey)

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

	return gjson.Get(string(body), "orders").String(), nil
}

func GetReservoirCollectionDaily(logger *golog.Logger, contract string) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://api.reservoir.tools/collections/daily-volumes/v1?id=" + contract
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().Reservoir.ApiKey)

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

	return gjson.Get(string(body), "collections").String(), nil
}
