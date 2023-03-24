package api

import (
	"github.com/kataras/iris/v12"
)

func GetNFTInfo(ctx iris.Context) {
	logger := ctx.Application().Logger()
	contract := ctx.Params().Get("contract")
	tokenID := ctx.Params().Get("token_id")

	data := GetNFTScan("https://restapi.nftscan.com/api/v2/assets/"+contract+"/"+tokenID+"?show_attribute=true", logger)

	ctx.WriteString(data)
}
