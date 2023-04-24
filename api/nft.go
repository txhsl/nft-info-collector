package api

import (
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
	offset, err := ctx.Params().GetInt("offset")
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	limit, err := ctx.Params().GetInt("limit")
	if err != nil {
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
	ctx.WriteString(data)
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
	ctx.WriteString(data)
}

func SearchNFTs(ctx iris.Context) {

}
