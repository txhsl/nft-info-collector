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
	collectionsAPI := app.Party("/collections")
	{
		collectionsAPI.Get("/immediate", api.ListImmediateCollections)
		collectionsAPI.Get("/cached", api.ListCachedCollections)
	}
	collectionAPI := app.Party("/collection")
	{
		collectionAPI.Get("/info/{contract}", api.GetCollectionInfo)
		collectionAPI.Get("/info/{contract}/{token_id}", api.GetNFTInfo)
		collectionAPI.Get("/statistics/{contract}", api.GetCollectionStatistics)
	}
	userAPI := app.Party("/user")
	{
		userAPI.Get("/nfts/{address}", api.GetUserNFTs)
		userAPI.Get("/collections/{address}/{page}", api.GetUserCollections)
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
