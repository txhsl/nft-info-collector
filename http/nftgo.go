package http

import (
	"io"
	"net/http"
	"nft-info-collector/config"

	"github.com/kataras/golog"
	"github.com/tidwall/gjson"
)

// Only used in DB cache
func GetNFTGoCollections(logger *golog.Logger, offset string, limit string) string {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/market/rank/collection/7d?by=volume&with_rarity=false&asc=false&offset=" + offset + "&limit=" + limit
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

// Used in immediate response
func GetNFTGoCollectionMetrics(logger *golog.Logger, contract string) string {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/collection/" + contract + "/metrics"
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
	return string(body)
}

// Used in immediate response
func GetNFTGoCollectionNFTs(logger *golog.Logger, contract string, offset string, limit string) string {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/collection/" + contract + "/nfts?offset=" + offset + "&limit=" + limit
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
	return gjson.Get(string(body), "nfts").String()
}

// Used in immediate response
func GetNFTGoNFTRarity(logger *golog.Logger, contract string, id string) string {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/nft/" + contract + "/" + id + "/rarity"
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
	return string(body)
}

// Used in immediate response
func GetNFTGoNFTMetrics(logger *golog.Logger, contract string, id string) string {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/nft/" + contract + "/" + id + "/metrics"
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
	return string(body)
}

// Used in immediate response
func GetNFTGoUserAssets(logger *golog.Logger, account string, offset string, limit string) string {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/address/" + account + "/portfolio?offset=" + offset + "&limit=" + limit
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
	return gjson.Get(string(body), "assets").String()
}
