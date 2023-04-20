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

	// fetch data
	data, err := http.GetNFTScanTrends(logger)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch trends")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}

	// deserialize result
	var collections []interface{}
	err = json.Unmarshal([]byte(data), &collections)
	if err != nil {
		logger.Error("[API] Failed to serialize collections")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}

	// connect db
	dbClient, err := db.Connect()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}
	defer dbClient.Disconnect(context.TODO())

	// cache result
	coll := dbClient.Database("nft-info-collector").Collection("trends")
	err = db.ReplaceCachedCollections(context.TODO(), logger, coll, collections)
	if err != nil {
		logger.Error("[DB] Failed to replace cached trends")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}

	ctx.WriteString(data)
}

func ListCachedTrends(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// connect db
	client, err := db.Connect()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}
	defer client.Disconnect(context.TODO())

	// search db
	coll := client.Database("nft-info-collector").Collection("trends")
	collections, err := db.GetCachedCollections(context.TODO(), logger, coll, 0, 100)
	if err != nil {
		logger.Error("[API] Failed to read cached trends")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}

	// serialize collections
	result, err := json.Marshal(collections)
	if err != nil {
		logger.Error("[API] Failed to deserialize cached trends")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}
	ctx.WriteString(string(result))
}
