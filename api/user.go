package api

import (
	"encoding/json"
	"nft-info-collector/http"

	"github.com/kataras/iris/v12"
)

func GetUserNFTs(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	address := ctx.Params().GetString("address")
	cursor := ctx.URLParam("cursor")
	if address == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// fetch data
	data, err := http.GetOpenSeaUserAssets(logger, address, cursor)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch user nfts")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	var assets map[string]interface{}
	err = json.Unmarshal([]byte(data), &assets)
	if err != nil {
		logger.Error("[API] Failed to deserialize user nfts")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}

	ctx.JSON(assets)
}
