package api

import (
	"nft-info-collector/http"
	"strconv"

	"github.com/kataras/iris/v12"
)

func GetNFTInfo(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	contract := ctx.Params().Get("contract")
	if contract == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
	}
	tokenID, err := strconv.Atoi(ctx.Params().Get("token_id"))
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
	}

	// fetch data
	data, err := http.GetOpenSeaAsset(logger, contract, tokenID)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch nft info")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}
	ctx.WriteString(data)
}
