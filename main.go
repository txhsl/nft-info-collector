package main

import (
	"nft-info-collector/api"
	"nft-info-collector/config"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	// iris config
	app := iris.New()
	app.UseRouter(recover.New())
	customLogger := logger.New(logger.Config{
		Status: true,
		Method: true,
		Path:   true,
	})
	app.Use(customLogger)

	// iris routes
	app.Get("/", hello)
	trendsAPI := app.Party("/trends")
	{
		trendsAPI.Get("/immediate", api.ListImmediateTrends)
		trendsAPI.Get("/cached", api.ListCachedTrends)
	}
	collectionAPI := app.Party("/collection")
	{
		collectionAPI.Get("/update", api.UpdateCachedCollections)
		collectionAPI.Get("/list", api.ListCachedCollections)
		collectionAPI.Get("/info/{contract}", api.GetCollectionInfo)
		collectionAPI.Get("/info/{contract}/{token_id}", api.GetNFTInfo)
	}
	userAPI := app.Party("/user")
	{
		userAPI.Get("/assets/{address}", api.GetUserNFTs)
	}
	searchAPI := app.Party("/search")
	{
		searchAPI.Get("/nfts/{keyword}", api.SearchNFTs)
		searchAPI.Get("/collections/{keyword}", api.SearchCollections)
	}
	app.Listen(config.Load().Port)
}

func hello(ctx iris.Context) {
	ctx.WriteString("Info collector is working.")
}
