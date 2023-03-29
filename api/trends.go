package api

import (
	"context"
	"encoding/json"
	"nft-info-collector/db"
	"nft-info-collector/http"

	"github.com/kataras/iris/v12"
)

func ListImmediateTrends(ctx iris.Context) {
	logger := ctx.Application().Logger()

	data := http.GetNFTScanTrends(logger)

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
	coll := dbClient.Database("nft-info-collector").Collection("trends")
	err = db.ReplaceCachedCollections(context.TODO(), logger, coll, collections)
	if err != nil {
		logger.Error("[DB] Failed to replace cached trends")
		panic(err)
	}

	ctx.WriteString(data)
}

func ListCachedTrends(ctx iris.Context) {
	logger := ctx.Application().Logger()
	// search db
	client, err := db.Connect()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		panic(err)
	}
	defer client.Disconnect(context.TODO())
	coll := client.Database("nft-info-collector").Collection("trends")
	collections, err := db.GetCachedCollections(context.TODO(), logger, coll)
	if err != nil {
		logger.Error("[API] Failed to read cached trends")
		panic(err)
	}
	// serialize collections
	result, err := json.Marshal(collections)
	if err != nil {
		logger.Error("[API] Failed to deserialize cached trends")
		panic(err)
	}
	ctx.WriteString(string(result))
}
