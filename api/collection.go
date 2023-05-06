package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"nft-info-collector/config"
	"nft-info-collector/db"
	"nft-info-collector/http"
	"strconv"
	"time"

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

			// type format
			info["dev_seller_fee_basis_points"], err = strconv.Atoi(info["dev_seller_fee_basis_points"].(string))
			if err != nil {
				logger.Error("[HTTP] Failed to format dev_seller_fee_basis_points")
				ctx.StopWithStatus(iris.StatusInternalServerError)
				return
			}
			info["dev_buyer_fee_basis_points"], err = strconv.Atoi(info["dev_buyer_fee_basis_points"].(string))
			if err != nil {
				logger.Error("[HTTP] Failed to format dev_buyer_fee_basis_points")
				ctx.StopWithStatus(iris.StatusInternalServerError)
				return
			}
			// info["opensea_seller_fee_basis_points"], err = strconv.Atoi(info["opensea_seller_fee_basis_points"].(string))
			// if err != nil {
			// 	logger.Error("[HTTP] Failed to format opensea_seller_fee_basis_points")
			// 	ctx.StopWithStatus(iris.StatusInternalServerError)
			// 	return
			// }
			info["opensea_buyer_fee_basis_points"], err = strconv.Atoi(info["opensea_buyer_fee_basis_points"].(string))
			if err != nil {
				logger.Error("[HTTP] Failed to format opensea_buyer_fee_basis_points")
				ctx.StopWithStatus(iris.StatusInternalServerError)
				return
			}
			createdDate, err := time.Parse(time.RFC3339, info["created_date"].(string))
			if err != nil {
				logger.Error("[HTTP] Failed to format created_date")
				ctx.StopWithStatus(iris.StatusInternalServerError)
				return
			}
			info["created_date"] = createdDate.Unix()
			info["last_updated"] = time.Now().Unix()
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

// TODO: update and read from db
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

// TODO: update and read from db
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

// TODO: update and read from db
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
	logger := ctx.Application().Logger()
	conf := config.Load().Keywords

	// filter params
	keyword := ctx.URLParam("keyword")
	timeRange := ctx.URLParam("time_range")
	isBadReq := true
	for _, t := range conf.Times {
		if t == timeRange {
			isBadReq = false
			break
		}
	}
	if isBadReq {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// range params
	floorPriceMin := ctx.URLParamIntDefault("floor_price_min", 0)
	floorPriceMax := ctx.URLParamIntDefault("floor_price_max", math.MaxInt)
	saleCountMin := ctx.URLParamIntDefault("sale_count_min", 0)
	saleCountMax := ctx.URLParamIntDefault("sale_count_max", math.MaxInt)
	royaltyMin := ctx.URLParamIntDefault("royalty_min", 0)
	royaltyMax := ctx.URLParamIntDefault("royalty_max", 10000)
	profitMarginMin := ctx.URLParamIntDefault("profit_margin_min", math.MinInt)
	profitMarginMax := ctx.URLParamIntDefault("profit_margin_max", math.MaxInt)
	ownerPercentageMin := ctx.URLParamIntDefault("owner_percentage_min", 0)
	ownerPercentageMax := ctx.URLParamIntDefault("owner_percentage_max", 100)
	collectionAgeMin := ctx.URLParamIntDefault("collection_age_min", 0)
	collectionAgeMax := ctx.URLParamIntDefault("collection_age_max", math.MaxInt)

	// sort params
	sort := ctx.URLParamDefault("sort", "total_volume")
	isBadReq = true
	for _, s := range conf.Sorts {
		if s == sort {
			isBadReq = false
			break
		}
	}
	if isBadReq {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	asc := ctx.URLParamBoolDefault("asc", false)

	// pagination params
	offset := ctx.URLParamIntDefault("offset", 0)
	limit := ctx.URLParamIntDefault("limit", 20)

	// search
	collections, err := db.SearchCollections(
		context.TODO(),
		logger,
		timeRange,
		floorPriceMin,
		floorPriceMax,
		saleCountMin,
		saleCountMax,
		royaltyMin,
		royaltyMax,
		profitMarginMin,
		profitMarginMax,
		ownerPercentageMin,
		ownerPercentageMax,
		collectionAgeMin,
		collectionAgeMax,
		keyword,
		sort,
		asc,
		offset,
		limit,
	)
	if err != nil {
		logger.Error(err)
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
