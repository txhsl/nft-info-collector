package http

import (
	"fmt"
	"io"
	"net/http"
	"nft-info-collector/config"

	"github.com/kataras/golog"
	"github.com/tidwall/gjson"
)

// Only used in DB cache
func GetNFTGoCollections(logger *golog.Logger, timeRange string, offset int, limit int) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/market/rank/collection/" + timeRange + "?offset=" + fmt.Sprint(offset) + "&limit=" + fmt.Sprint(limit)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().NFTGo.ApiKey)

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

// @deprecated
func GetNFTGoCollectionInfo(logger *golog.Logger, contract string) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/collection/" + contract + "/info"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().NFTGo.ApiKey)

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
	return string(body), nil
}

// Used in immediate response
// TODO: only use in DB cache
func GetNFTGoCollectionMetrics(logger *golog.Logger, contract string) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/collection/" + contract + "/metrics"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().NFTGo.ApiKey)

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
	return string(body), nil
}

// Used in immediate response
func GetNFTGoCollectionNFTs(logger *golog.Logger, contract string, offset int, limit int) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/collection/" + contract + "/nfts?offset=" + fmt.Sprint(offset) + "&limit=" + fmt.Sprint(limit)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().NFTGo.ApiKey)

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
	return gjson.Get(string(body), "nfts").String(), nil
}

// Used in immediate response
func GetNFTGoNFTRarity(logger *golog.Logger, contract string, id int) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/nft/" + contract + "/" + fmt.Sprint(id) + "/rarity"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().NFTGo.ApiKey)

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
	return string(body), nil
}

// Used in immediate response
func GetNFTGoNFTMetrics(logger *golog.Logger, contract string, id int) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/nft/" + contract + "/" + fmt.Sprint(id) + "/metrics"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().NFTGo.ApiKey)

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
	return string(body), nil
}

// Used in immediate response
func GetNFTGoUserAssets(logger *golog.Logger, account string, offset int, limit int) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://data-api.nftgo.io/eth/v1/address/" + account + "/portfolio?offset=" + fmt.Sprint(offset) + "&limit=" + fmt.Sprint(limit)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().NFTGo.ApiKey)

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
	return gjson.Get(string(body), "assets").String(), nil
}
