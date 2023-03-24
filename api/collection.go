package api

import (
	"context"
	"encoding/json"
	"nft-info-collector/db"

	"github.com/kataras/iris/v12"
)

func ListImmediateCollections(ctx iris.Context) {
	logger := ctx.Application().Logger()

	data := GetNFTScan("https://restapi.nftscan.com/api/v2/statistics/ranking/trade?time=7d&sort_field=volume&sort_direction=desc&show_7d_trends=false", logger)

	// cache result
	var collections []interface{}
	err := json.Unmarshal([]byte(data), &collections)
	if err != nil {
		logger.Error("[API] Failed to serialize collections")
		panic(err)
	}
	dbClient, err := db.Connect()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		panic(err)
	}
	defer dbClient.Disconnect(context.TODO())
	coll := dbClient.Database("nft-info-collector").Collection("collections")
	db.CacheCollections(context.TODO(), logger, coll, collections)

	ctx.WriteString(data)
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
	contract := ctx.Params().Get("contract")

	data := GetNFTScan("https://restapi.nftscan.com/api/v2/collections/"+contract+"?show_attribute=true", logger)

	ctx.WriteString(data)
}

func GetCollectionStatistics(ctx iris.Context) {
	logger := ctx.Application().Logger()
	contract := ctx.Params().Get("contract")

	data := GetNFTScan("https://restapi.nftscan.com/api/v2/statistics/collection/"+contract, logger)

	ctx.WriteString(data)
}
