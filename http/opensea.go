package http

import (
	"fmt"
	"io"
	"net/http"
	"nft-info-collector/config"

	"github.com/kataras/golog"
)

// Only used in DB cache
// @deprecated
func GetOpenSeaCollections(logger *golog.Logger, offset int, limit int) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://api.opensea.io/api/v1/collections?format=json&offset=" + fmt.Sprint(offset) + "&limit=" + fmt.Sprint(limit)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")

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

// Only used in DB cache
func GetOpenSeaCollectionInfo(logger *golog.Logger, slug string) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://api.opensea.io/api/v1/collection/" + slug + "?format=json"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

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

func GetOpenSeaCollectionOffers(logger *golog.Logger, slug string) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://api.opensea.io/v2/offers/collection/" + slug + "?format=json"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().OpenSea.ApiKey)

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

func GetOpenSeaCollectionRecentSales(logger *golog.Logger, slug string) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://api.opensea.io/api/v1/events?collection_slug=" + slug + "&event_type=successful&only_opensea=true"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().OpenSea.ApiKey)

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
func GetOpenSeaAsset(logger *golog.Logger, contract string, id int) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://api.opensea.io/api/v1/asset/" + contract + "/" + fmt.Sprint(id) + "/?format=json&include_orders=false"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().OpenSea.ApiKey)

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
func GetOpenSeaUserAssets(logger *golog.Logger, account string, cursor string) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://api.opensea.io/api/v1/assets?format=json&owner=" + account
	if cursor != "" {
		url += "&cursor=" + cursor
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-API-KEY", config.Load().OpenSea.ApiKey)

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
