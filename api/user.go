package api

import (
	"nft-info-collector/http"

	"github.com/kataras/iris/v12"
)

func GetUserNFTs(ctx iris.Context) {
	logger := ctx.Application().Logger()
	account := ctx.Params().Get("account")

	data := http.GetOpenSeaAssets(logger, account)

	ctx.WriteString(data)
}
