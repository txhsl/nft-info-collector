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

	// parse params
	timeRange := ctx.Params().GetString("time_range")
	isBadReq := true
	for _, t := range conf.TimeRanges {
		if t == timeRange {
			isBadReq = false
			break
		}
	}
	if isBadReq {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// connect db
	dbClient, err := db.Connect()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	defer dbClient.Disconnect(context.TODO())
	coll := dbClient.Database("nft-info-collector").Collection("collections-" + timeRange)

	for i := 0; i < conf.Limit; i += conf.PageSize {
		// fetch data
		data, err := http.GetNFTGoCollections(logger, timeRange, i, conf.PageSize)
		if err != nil {
			logger.Error("[HTTP] Failed to fetch collections")
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

		// cache result
		err = db.UpdateCachedCollections(context.TODO(), logger, coll, collections)
		if err != nil {
			logger.Error("[DB] Failed to update cached collections")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}

		logger.Info("[DB] Collections updated: " + fmt.Sprint(i+conf.PageSize) + " / " + fmt.Sprint(conf.Limit))
	}

	ctx.WriteString("OK")
}

func ListCachedCollections(ctx iris.Context) {
	logger := ctx.Application().Logger()
	conf := config.Load().NFTGo

	// parse params
	timeRange := ctx.Params().GetString("time_range")
	isBadReq := true
	for _, t := range conf.TimeRanges {
		if t == timeRange {
			isBadReq = false
			break
		}
	}
	if isBadReq {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	keyword := ctx.Params().GetString("keyword")
	isBadReq = true
	for _, k := range conf.Keywords {
		if k == keyword {
			isBadReq = false
			break
		}
	}
	if isBadReq {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	asc, err := ctx.Params().GetBool("asc")
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	offset, err := ctx.Params().GetInt("offset")
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	limit, err := ctx.Params().GetInt("limit")
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// search db
	collections, err := db.GetSortedCollections(context.TODO(), logger, timeRange, keyword, asc, offset, limit)
	if err != nil {
		logger.Error("[DB] Failed to read cached collections")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	// serialize result
	result, err := json.Marshal(collections)
	if err != nil {
		logger.Error("[API] Failed to deserialize cached collections")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	ctx.WriteString(string(result))
}

func FilterCachedCollections(ctx iris.Context) {
	logger := ctx.Application().Logger()
	conf := config.Load().NFTGo

	// parse params
	timeRange := ctx.Params().GetString("time_range")
	isBadReq := true
	for _, t := range conf.TimeRanges {
		if t == timeRange {
			isBadReq = false
			break
		}
	}
	if isBadReq {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	filter := ctx.Params().GetString("filter")
	isBadReq = true
	for _, f := range conf.Filters {
		if f == filter {
			isBadReq = false
			break
		}
	}
	if isBadReq {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	value := ctx.Params().GetString("value")
	if value == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	offset, err := ctx.Params().GetInt("offset")
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	limit, err := ctx.Params().GetInt("limit")
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// search db
	collections, err := db.GetFilteredCollections(context.TODO(), logger, timeRange, filter, value, offset, limit)
	if err != nil {
		logger.Error("[DB] Failed to read cached collections")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	// serialize result
	result, err := json.Marshal(collections)
	if err != nil {
		logger.Error("[API] Failed to deserialize cached collections")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	ctx.WriteString(string(result))
}

func GetCollectionDetail(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	slug := ctx.Params().GetString("slug")
	if slug == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// fetch data
	data, err := http.GetOpenSeaCollectionInfo(logger, slug)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch collection info")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	ctx.WriteString(data)
}

func GetCollectionInfo(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	contract := ctx.Params().GetString("contract")
	if contract == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// fetch data
	data, err := http.GetNFTGoCollectionInfo(logger, contract)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch collection metrics")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	ctx.WriteString(data)
}

func GetCollectionMetrics(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	contract := ctx.Params().GetString("contract")
	if contract == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// fetch data
	data, err := http.GetNFTGoCollectionMetrics(logger, contract)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch collection metrics")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	ctx.WriteString(data)
}

func GetCollectionNFTs(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	contract := ctx.Params().GetString("contract")
	if contract == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	offset, err := ctx.Params().GetInt("offset")
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	limit, err := ctx.Params().GetInt("limit")
	if err != nil {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// fetch data
	data, err := http.GetNFTGoCollectionNFTs(logger, contract, offset, limit)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch collection nfts")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	ctx.WriteString(data)
}
