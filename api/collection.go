package api

import (
	"context"
	"encoding/json"
	"fmt"
	"nft-info-collector/config"
	"nft-info-collector/db"
	"nft-info-collector/http"
	"strconv"

	"github.com/kataras/iris/v12"
)

func UpdateCachedCollections(ctx iris.Context) {
	logger := ctx.Application().Logger()
	conf := config.Load().NFTGo

	// connect db
	dbClient, err := db.Connect()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}
	defer dbClient.Disconnect(context.TODO())
	coll := dbClient.Database("nft-info-collector").Collection("collections")

	for i := 0; i < conf.Limit; i += conf.PageSize {
		// fetch data
		data, err := http.GetNFTGoCollections(logger, i, conf.PageSize)
		if err != nil {
			logger.Error("[HTTP] Failed to fetch collections")
			ctx.StopWithStatus(iris.StatusInternalServerError)
		}

		// deserialize result
		var collections []interface{}
		err = json.Unmarshal([]byte(data), &collections)
		if err != nil {
			logger.Error("[API] Failed to serialize collections")
			ctx.StopWithStatus(iris.StatusInternalServerError)
		}

		// cache result
		err = db.UpdateCachedCollections(context.TODO(), logger, coll, collections)
		if err != nil {
			logger.Error("[DB] Failed to update cached collections")
			ctx.StopWithStatus(iris.StatusInternalServerError)
		}

		logger.Info("[DB] Collections updated: " + fmt.Sprint(i+conf.PageSize) + " / " + fmt.Sprint(conf.Limit))
	}

	ctx.WriteString("OK")
}

func ListCachedCollections(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	offset, err := strconv.Atoi(ctx.Params().Get("offset"))
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
	}
	limit, err := strconv.Atoi(ctx.Params().Get("limit"))
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
	}

	// connect db
	client, err := db.Connect()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}
	defer client.Disconnect(context.TODO())

	// search db
	coll := client.Database("nft-info-collector").Collection("collections")
	collections, err := db.GetCachedCollections(context.TODO(), logger, coll, offset, limit)
	if err != nil {
		logger.Error("[DB] Failed to read cached collections")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}

	// serialize result
	result, err := json.Marshal(collections)
	if err != nil {
		logger.Error("[API] Failed to deserialize cached collections")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}
	ctx.WriteString(string(result))
}

func GetCollectionInfo(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	slug := ctx.Params().Get("slug")
	if slug == "" {
		ctx.StatusCode(iris.StatusBadRequest)
	}

	// fetch data
	data, err := http.GetOpenSeaCollectionInfo(logger, slug)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch collection info")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}
	ctx.WriteString(data)
}
