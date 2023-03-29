package api

import (
	"nft-info-collector/http"

	"github.com/kataras/iris/v12"
)

func GetNFTInfo(ctx iris.Context) {
	logger := ctx.Application().Logger()
	contract := ctx.Params().Get("contract")
	tokenID := ctx.Params().Get("token_id")

	data := http.GetOpenSeaAsset(logger, contract, tokenID)

	ctx.WriteString(data)
}
