package api

import (
	"encoding/json"
	"nft-info-collector/http"

	"github.com/kataras/iris/v12"
)

func GetCollectionNFTs(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	contract := ctx.Params().GetString("contract")
	if contract == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	offset := ctx.URLParamIntDefault("offset", 0)
	if offset < 0 {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	limit := ctx.URLParamIntDefault("limit", 20)
	if limit < 0 || limit > 50 {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// fetch data
	data, err := http.GetNFTGoCollectionNFTs(logger, contract, offset, limit)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch collection nfts")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	var nfts map[string]interface{}
	err = json.Unmarshal([]byte(data), &nfts)
	if err != nil {
		logger.Error("[API] Failed to deserialize collection nfts")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}
	ctx.JSON(map[string]interface{}{"nfts": nfts["nfts"]})
}

func GetNFTDetail(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	contract := ctx.Params().GetString("contract")
	if contract == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	tokenID, err := ctx.Params().GetInt("token_id")
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// fetch data
	data, err := http.GetOpenSeaAsset(logger, contract, tokenID)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch nft info")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	var asset map[string]interface{}
	err = json.Unmarshal([]byte(data), &asset)
	if err != nil {
		logger.Error("[API] Failed to deserialize nft info")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}
	ctx.JSON(asset)
}
