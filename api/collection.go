package api

import (
	"context"
	"encoding/json"
	"fmt"
	"nft-info-collector/config"
	"nft-info-collector/db"
	"nft-info-collector/http"

	"github.com/kataras/iris/v12"
)

func UpdateCachedCollections(ctx iris.Context) {
	logger := ctx.Application().Logger()
	conf := config.Load().NFTGo

	dbClient, err := db.Connect()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		panic(err)
	}
	defer dbClient.Disconnect(context.TODO())
	coll := dbClient.Database("nft-info-collector").Collection("collections")

	for i := 0; i < conf.Limit; i += conf.PageSize {
		data := http.GetNFTGoCollections(logger, i, conf.PageSize)

		// cache result
		var collections []interface{}
		err := json.Unmarshal([]byte(data), &collections)
		if err != nil {
			logger.Error("[API] Failed to serialize collections")
			panic(err)
		}
		err = db.UpdateCachedCollections(context.TODO(), logger, coll, collections)
		if err != nil {
			logger.Error("[DB] Failed to update cached collections")
			panic(err)
		}

		logger.Info("[DB] Collections updated: " + fmt.Sprint(i+conf.PageSize) + " / " + fmt.Sprint(conf.Limit))
	}

	ctx.WriteString("OK")
}

func ListCachedCollections(ctx iris.Context) {
	logger := ctx.Application().Logger()
	// search db
	client, err := db.Connect()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		panic(err)
	}
	defer client.Disconnect(context.TODO())
	coll := client.Database("nft-info-collector").Collection("collections")
	collections, err := db.GetCachedCollections(context.TODO(), logger, coll)
	if err != nil {
		logger.Error("[API] Failed to read cached collections")
		panic(err)
	}
	// serialize collections
	result, err := json.Marshal(collections)
	if err != nil {
		logger.Error("[API] Failed to deserialize cached collections")
		panic(err)
	}
	ctx.WriteString(string(result))
}

func GetCollectionInfo(ctx iris.Context) {
	logger := ctx.Application().Logger()
	slug := ctx.Params().Get("slug")

	data := http.GetOpenSeaCollectionInfo(logger, slug)

	ctx.WriteString(data)
}
