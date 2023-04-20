package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kataras/golog"
	"github.com/tidwall/gjson"
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
	return gjson.Get(string(body), "collections").String(), nil
}

// Used in immediate response
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
	return gjson.Get(string(body), "collection").String(), nil
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
func GetOpenSeaUserAssets(logger *golog.Logger, account string) (string, error) {
	// build request
	httpClient := &http.Client{}
	url := "https://api.opensea.io/api/v1/assets?format=json&owner=" + account
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
	assets := gjson.Get(string(body), "assets").Array()
	cursor := gjson.Get(string(body), "next")

	// get left pages
	for {
		if cursor.Value() == nil {
			break
		}
		url = "https://api.opensea.io/api/v1/assets?format=json&owner=" + account + "&cursor=" + cursor.String()
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return "", err
		}
		res, err := httpClient.Do(req)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		assets = append(assets, gjson.Get(string(body), "assets").Array()...)
	}

	// serialize
	data := "["
	for i := 0; i < len(assets); i++ {
		data += assets[i].String()
	}
	return data + "]", nil
}
