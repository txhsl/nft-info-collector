package api

import (
	"github.com/kataras/iris/v12"
)

func GetUserNFTs(ctx iris.Context) {
	logger := ctx.Application().Logger()
	address := ctx.Params().Get("address")

	data := GetNFTScan("https://restapi.nftscan.com/api/v2/account/own/all/"+address+"?erc_type=erc721&show_attribute=false", logger)

	ctx.WriteString(data)
}

func GetUserCollections(ctx iris.Context) {
	logger := ctx.Application().Logger()
	address := ctx.Params().Get("address")
	page := ctx.Params().Get("page")

	data := GetNFTScan("https://restapi.nftscan.com/api/v2/collections/own/"+address+"?erc_type=erc721&offest="+page, logger)

	ctx.WriteString(data)
}
