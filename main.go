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
	updateAPI := app.Party("/update")
	{
		updateAPI.Get("/index", api.UpdateCachedCollectionIndex)
		updateAPI.Get("/details", api.UpdateCachedCollectionDetails)
		updateAPI.Get("/metrics", api.UpdateCachedCollectionMetrics)
	}
	collectionAPI := app.Party("/collection")
	{
		collectionAPI.Get("/sort/{time_range}/{keyword}/{asc:boolean}/{offset:int}/{limit:int}", api.SortCachedCollections)
		collectionAPI.Get("/filter/{time_range}/{filter}/{value}/{offset:int}/{limit:int}", api.FilterCachedCollections)
		collectionAPI.Get("/search/{keyword}/{min}/{max}", api.SearchCollections)

		collectionAPI.Get("/info/{contract}", api.GetCollectionInfo)
		collectionAPI.Get("/metrics/{contract}", api.GetCollectionMetrics)
		collectionAPI.Get("/detail/{slug}", api.GetCollectionDetail)
	}
	nftAPI := app.Party("/nft")
	{
		nftAPI.Get("/list/{contract}/{offset:int}/{limit:int}", api.GetCollectionNFTs)
		nftAPI.Get("/search/{contract}/{keyword}/{min}/{max}", api.SearchNFTs)
		// Banned by Cloudflare with 1020 - Access denied
		nftAPI.Get("/detail/{contract}/{token_id}", api.GetNFTDetail)
	}
	userAPI := app.Party("/user")
	{
		userAPI.Get("/assets/{address}", api.GetUserNFTs)
	}
	app.Listen(config.Load().Port)
}

func hello(ctx iris.Context) {
	ctx.WriteString("Info collector is working.")
}
