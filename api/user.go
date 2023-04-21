package api

import (
	"nft-info-collector/http"

	"github.com/kataras/iris/v12"
)

func GetUserNFTs(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	address := ctx.Params().GetString("address")
	if address == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// fetch data
	data, err := http.GetOpenSeaUserAssets(logger, address)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch user nfts")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	ctx.WriteString(data)
}
