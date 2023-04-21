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
		collectionAPI.Get("/update/{time_range}", api.UpdateCachedCollections)
		collectionAPI.Get("/list/{time_range}/{keyword}/{asc:boolean}/{offset:int}/{limit:int}", api.ListCachedCollections)
		collectionAPI.Get("/filter/{time_range}/{filter}/{value}/{offset:int}/{limit:int}", api.FilterCachedCollections)
		collectionAPI.Get("/info/{contract}", api.GetCollectionInfo)
		collectionAPI.Get("/metrics/{contract}", api.GetCollectionMetrics)
		collectionAPI.Get("/detail/{slug}", api.GetCollectionDetail)
		collectionAPI.Get("/nfts/{contract}/{offset:int}/{limit:int}", api.GetCollectionNFTs)
	}
	nftAPI := app.Party("/nft")
	{
		// Banned by Cloudflare with 1020 - Access denied
		nftAPI.Get("/detail/{contract}/{token_id}", api.GetNFTDetail)
	}
	userAPI := app.Party("/user")
	{
		userAPI.Get("/assets/{address}", api.GetUserNFTs)
	}
	searchAPI := app.Party("/search")
	{
		searchAPI.Get("/nfts/{contract}/{keyword}/{min}/{max}", api.SearchNFTs)
		searchAPI.Get("/collections/{keyword}/{min}/{max}", api.SearchCollections)
	}
	app.Listen(config.Load().Port)
}

func hello(ctx iris.Context) {
	ctx.WriteString("Info collector is working.")
}
