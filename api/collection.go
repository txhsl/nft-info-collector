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

// DB related, write
func UpdateCachedCollectionIndex(ctx iris.Context) {
	logger := ctx.Application().Logger()
	conf := config.Load().NFTGo

	// connect db
	dbClient, err := db.Connect()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	defer dbClient.Disconnect(context.TODO())
	coll := dbClient.Database("nft-info-collector").Collection("collection-index")

	for i := 0; i < conf.Limit; i += conf.PageSize {
		// fetch data
		data, err := http.GetNFTGoCollections(logger, "all", i, conf.PageSize)
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
		err = db.UpdateCollectionIndex(context.TODO(), logger, coll, collections)
		if err != nil {
			logger.Error("[DB] Failed to update cached collections")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}

		logger.Info("[DB] Index updated: " + fmt.Sprint(i+conf.PageSize) + " / " + fmt.Sprint(conf.Limit))
	}

	ctx.WriteString("OK")
}

// DB related, write and read
func UpdateCachedCollectionMetrics(ctx iris.Context) {
	logger := ctx.Application().Logger()
	conf := config.Load().Reservoir

	// connect db
	dbClient, err := db.Connect()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	defer dbClient.Disconnect(context.TODO())

	// fetch details
	coll := dbClient.Database("nft-info-collector").Collection("collection-metrics")

	for i := 0; i < conf.Limit; i += conf.PageSize {
		// read db
		result, err := db.GetSortedCollectionIndex(context.TODO(), logger, "volume_usd", false, i, conf.PageSize)
		if err != nil {
			logger.Error("[DB] Failed to read cached collections")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(result)
		if err != nil {
			logger.Error("[API] Failed to deserialize cached collections")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}
		var batch []map[string]interface{}
		err = json.Unmarshal([]byte(data), &batch)
		if err != nil {
			logger.Error("[API] Failed to get collection contracts")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}

		contracts := make([]string, 0, len(batch))
		for _, collection := range batch {
			contract := collection["contracts"].([]interface{})[0].(string)
			contracts = append(contracts, contract)
		}
		collections, err := http.GetReservoirCollections(logger, contracts)
		if err != nil {
			logger.Error("[HTTP] Failed to fetch collection info")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}

		// deserialize result
		var metrics []interface{}
		err = json.Unmarshal([]byte(collections), &metrics)
		if err != nil {
			logger.Error("[API] Failed to serialize collections")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}

		// cache result
		err = db.UpdateCollectionMetrics(context.TODO(), logger, coll, metrics)
		if err != nil {
			logger.Error("[DB] Failed to update cached collection details")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}

		logger.Info("[DB] Metrics updated: " + fmt.Sprint(i+conf.PageSize) + " / " + fmt.Sprint(conf.Limit))
	}
}

// DB related, write and read
func UpdateCachedCollectionDetails(ctx iris.Context) {
	logger := ctx.Application().Logger()
	conf := config.Load().OpenSea

	// connect db
	dbClient, err := db.Connect()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	defer dbClient.Disconnect(context.TODO())

	// fetch details
	coll := dbClient.Database("nft-info-collector").Collection("collection-details")

	for i := 0; i < conf.Limit; i += conf.PageSize {
		// read db
		result, err := db.GetSortedCollectionIndex(context.TODO(), logger, "volume_usd", false, i, conf.PageSize)
		if err != nil {
			logger.Error("[DB] Failed to read cached collections")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(result)
		if err != nil {
			logger.Error("[API] Failed to deserialize cached collections")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}
		var collections []map[string]interface{}
		err = json.Unmarshal([]byte(data), &collections)
		if err != nil {
			logger.Error("[API] Failed to get collection slugs")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}

		batch := []map[string]interface{}{}
		for _, collection := range collections {
			// fetch data
			slug := collection["opensea_slug"].(string)
			data, err := http.GetOpenSeaCollectionInfo(logger, slug)
			if err != nil {
				logger.Error("[HTTP] Failed to fetch collection info")
				ctx.StopWithStatus(iris.StatusInternalServerError)
				return
			}

			// deserialize result
			var info map[string]interface{}
			err = json.Unmarshal([]byte(data), &info)
			if err != nil {
				// skip if not found
				continue
			}
			batch = append(batch, info)
		}

		// cache result
		err = db.UpdateCollectionDetails(context.TODO(), logger, coll, batch)
		if err != nil {
			logger.Error("[DB] Failed to update cached collection details")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}

		logger.Info("[DB] Details updated: " + fmt.Sprint(i+conf.PageSize) + " / " + fmt.Sprint(conf.Limit))
	}
}

// DB related, read
// TODO: update
func SortCachedCollections(ctx iris.Context) {
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
	collections, err := db.GetSortedCollectionIndex(context.TODO(), logger, keyword, asc, offset, limit)
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

// DB related, read
// TODO: update
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
	collections, err := db.GetFilteredCollectionIndex(context.TODO(), logger, filter, value, offset, limit)
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

// TODO: read from db
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

// TODO: read from db
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

// TODO: update
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

func SearchCollections(ctx iris.Context) {

}
