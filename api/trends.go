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
		return
	}

	// deserialize result
	var collections []interface{}
	err = json.Unmarshal([]byte(data), &collections)
	if err != nil {
		logger.Error("[API] Failed to serialize collections")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	err = db.ReplaceCachedTrends(context.TODO(), logger, collections)
	if err != nil {
		logger.Error("[DB] Failed to replace cached trends")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	ctx.WriteString(data)
}

func ListCachedTrends(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// search db
	collections, err := db.GetCachedTrends(context.TODO(), logger, 0, 100)
	if err != nil {
		logger.Error("[DB] Failed to read cached trends")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	// serialize collections
	result, err := json.Marshal(collections)
	if err != nil {
		logger.Error("[API] Failed to deserialize cached trends")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	ctx.WriteString(string(result))
}
