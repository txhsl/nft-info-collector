package api

import (
	"context"
	"encoding/json"
	"fmt"
	"nft-info-collector/config"
	"nft-info-collector/db"
	"nft-info-collector/http"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/tidwall/gjson"
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
		err = json.Unmarshal([]byte(gjson.Get(data, "collections").String()), &collections)
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
		err = json.Unmarshal([]byte(gjson.Get(collections, "collections").String()), &metrics)
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
	ctx.WriteString("OK")
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

		detailBatch := []map[string]interface{}{}
		for _, collection := range collections {
			// fetch detail
			slug := collection["opensea_slug"].(string)
			data, err := http.GetOpenSeaCollectionInfo(logger, slug)
			if err != nil {
				logger.Error("[HTTP] Failed to fetch collection info")
				ctx.StopWithStatus(iris.StatusInternalServerError)
				return
			}

			// deserialize result
			var info map[string]interface{}
			err = json.Unmarshal([]byte(gjson.Get(data, "collection").String()), &info)
			if err != nil {
				// skip if not found
				continue
			}

			// read last update time
			lastUpdated, _ := db.GetDetailLastUpdated(context.TODO(), coll, slug)
			oneDaySales := info["stats"].(map[string]interface{})["one_day_sales"].(float64)
			if oneDaySales == 0 && lastUpdated >= time.Now().Add(-24*time.Hour).Unix() {
				// skip if no sales in 24 hours
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
			info["opensea_seller_fee_basis_points"] = int(info["opensea_seller_fee_basis_points"].(float64))
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

			// add total royalty
			if info["is_creator_fees_enforced"].(bool) {
				info["total_royalty"] = info["opensea_seller_fee_basis_points"].(int) + info["dev_seller_fee_basis_points"].(int)
			} else {
				info["total_royalty"] = info["opensea_seller_fee_basis_points"].(int)
			}

			// fetch collection offers if enabled
			info["top_bid_price"] = nil
			if info["is_collection_offers_enabled"].(bool) {
				data, err := http.GetOpenSeaCollectionOffers(logger, slug)
				if err != nil {
					logger.Error("[HTTP] Failed to fetch collection offers")
					ctx.StopWithStatus(iris.StatusInternalServerError)
					return
				}

				// find top bid
				for _, offer := range gjson.Get(data, "offers").Array() {
					chain := gjson.Get(offer.String(), "chain").String()
					token := gjson.Get(offer.String(), "protocol_data.parameters.offer.0.token").String()
					if chain == "ethereum" && token == "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" {
						info["top_bid_price"] = gjson.Get(offer.String(), "protocol_data.parameters.offer.0.endAmount").Float() / 1000000000000000000
					}
					break
				}
			}

			// add last updated time
			info["last_updated"] = time.Now().Unix()
			detailBatch = append(detailBatch, info)
		}

		// cache result
		err = db.UpdateCollectionDetails(context.TODO(), logger, coll, detailBatch)
		if err != nil {
			logger.Error("[DB] Failed to update cached collection details")
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}

		logger.Info("[DB] Details updated: " + fmt.Sprint(i+conf.PageSize) + " / " + fmt.Sprint(conf.Limit))
	}
	ctx.WriteString("OK")
}
