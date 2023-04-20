package api

import (
	"nft-info-collector/http"

	"github.com/kataras/iris/v12"
)

func GetUserNFTs(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	account := ctx.Params().Get("account")
	if account == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
	}

	// fetch data
	data, err := http.GetOpenSeaUserAssets(logger, account)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch user nfts")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}

	ctx.WriteString(data)
}
