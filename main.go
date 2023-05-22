package main

import (
	"nft-info-collector/api"
	"nft-info-collector/config"

	"github.com/iris-contrib/middleware/cors"
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

	// cors config
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	app.UseRouter(crs)

	// iris routes
	app.Get("/", hello)

	trendsAPI := app.Party("/trends")
	{
		trendsAPI.Get("/immediate", api.ListImmediateTrends)
		trendsAPI.Get("/cached", api.ListCachedTrends)
	}
	updateAPI := app.Party("/update")
	{
		updateAPI.Post("/index", api.UpdateCachedCollectionIndex)
		updateAPI.Post("/details", api.UpdateCachedCollectionDetails)
		updateAPI.Post("/metrics", api.UpdateCachedCollectionMetrics)
	}
	collectionAPI := app.Party("/collection")
	{
		collectionAPI.Post("/search", api.SearchCollections)
		collectionAPI.Get("/detail/{slug}", api.GetCollectionDetail)
	}
	nftAPI := app.Party("/nft")
	{
		nftAPI.Get("/list/{contract}", api.GetCollectionNFTs)
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
