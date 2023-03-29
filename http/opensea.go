package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kataras/golog"
	"github.com/tidwall/gjson"
)

func GetOpenSeaCollections(logger *golog.Logger, offset int, limit int) string {
	// build request
	httpClient := &http.Client{}
	url := "https://api.opensea.io/api/v1/collections?format=json&offset=" + fmt.Sprint(offset) + "&limit=" + fmt.Sprint(limit)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error("[API] Failed to build opensea request")
		panic(err)
	}
	req.Header.Add("Accept", "application/json")

	// send request
	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error("[API] Failed to send opensea request")
		panic(err)
	}
	defer res.Body.Close()

	// analysis response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("[API] Failed to read opensea response")
		panic(err)
	}
	return gjson.Get(string(body), "collections").String()
}

func GetOpenSeaCollectionInfo(logger *golog.Logger, slug string) string {
	// build request
	httpClient := &http.Client{}
	url := "https://api.opensea.io/api/v1/collection/" + slug + "?format=json"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error("[API] Failed to build opensea request")
		panic(err)
	}

	// send request
	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error("[API] Failed to send opensea request")
		panic(err)
	}
	defer res.Body.Close()

	// analysis response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("[API] Failed to read opensea response")
		panic(err)
	}
	return gjson.Get(string(body), "collection").String()
}

func GetOpenSeaAsset(logger *golog.Logger, contract string, id string) string {
	// build request
	httpClient := &http.Client{}
	url := "https://api.opensea.io/api/v1/asset/" + contract + "/" + id + "/?format=json&include_orders=false"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error("[API] Failed to build opensea request")
		panic(err)
	}

	// send request
	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error("[API] Failed to send opensea request")
		panic(err)
	}
	defer res.Body.Close()

	// analysis response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("[API] Failed to read opensea response")
		panic(err)
	}
	return string(body)
}

func GetOpenSeaAssets(logger *golog.Logger, account string) string {
	// build request
	httpClient := &http.Client{}
	url := "https://api.opensea.io/api/v1/assets?owner=" + account + "?format=json"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error("[API] Failed to build opensea request")
		panic(err)
	}

	// send request
	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error("[API] Failed to send opensea request")
		panic(err)
	}
	defer res.Body.Close()

	// analysis response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("[API] Failed to read opensea response")
		panic(err)
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
			logger.Error("[API] Failed to build opensea request")
			panic(err)
		}
		res, err := httpClient.Do(req)
		if err != nil {
			logger.Error("[API] Failed to send opensea request")
			panic(err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			logger.Error("[API] Failed to read opensea response")
			panic(err)
		}
		assets = append(assets, gjson.Get(string(body), "assets").Array()...)
	}

	// serialize
	data := "["
	for i := 0; i < len(assets); i++ {
		data += assets[i].String()
	}
	return data + "]"
}
