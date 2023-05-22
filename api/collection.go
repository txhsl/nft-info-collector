package api

import (
	"context"
	"encoding/json"
	"math"
	"nft-info-collector/config"
	"nft-info-collector/db"
	"nft-info-collector/http"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/tidwall/gjson"
)

type SearchRequest struct {
	Keyword            string  `json:"keyword"`
	TimeRange          string  `json:"time_range"`
	FloorPriceMin      float32 `json:"floor_price_min"`
	FloorPriceMax      float32 `json:"floor_price_max"`
	SaleCountMin       int     `json:"sale_count_min"`
	SaleCountMax       int     `json:"sale_count_max"`
	RoyaltyMin         int     `json:"royalty_min"`
	RoyaltyMax         int     `json:"royalty_max"`
	ProfitMarginMin    float32 `json:"profit_margin_min"`
	ProfitMarginMax    float32 `json:"profit_margin_max"`
	OwnerPercentageMin float32 `json:"owner_percentage_min"`
	OwnerPercentageMax float32 `json:"owner_percentage_max"`
	CollectionAgeMin   int     `json:"collection_age_min"`
	CollectionAgeMax   int     `json:"collection_age_max"`
	Sort               string  `json:"sort"`
	Asc                bool    `json:"asc"`
	Offset             int     `json:"offset"`
	Limit              int     `json:"limit"`
}

func SearchCollections(ctx iris.Context) {
	logger := ctx.Application().Logger()
	conf := config.Load().Keywords

	// parse params
	body, err := ctx.GetBody()
	if err != nil {
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	type searchReq SearchRequest
	req := &searchReq{
		Keyword:            "",
		TimeRange:          "1d",
		FloorPriceMin:      0.0,
		FloorPriceMax:      math.MaxFloat32,
		SaleCountMin:       0,
		SaleCountMax:       math.MaxInt,
		RoyaltyMin:         0,
		RoyaltyMax:         math.MaxInt,
		ProfitMarginMin:    -math.MaxFloat32,
		ProfitMarginMax:    math.MaxFloat32,
		OwnerPercentageMin: 0.0,
		OwnerPercentageMax: 100.0,
		CollectionAgeMin:   0,
		CollectionAgeMax:   math.MaxInt,
		Sort:               "total_volume",
		Asc:                false,
		Offset:             0,
		Limit:              20,
	}
	if err := json.Unmarshal(body, &req); err != nil {
		logger.Error("[API] Failed to deserialize search request")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	// validate params
	isBadReq := true
	for _, t := range conf.Times {
		if t == req.TimeRange {
			isBadReq = false
			break
		}
	}
	if isBadReq {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	isBadReq = true
	for _, s := range conf.Sorts {
		if s == req.Sort {
			isBadReq = false
			break
		}
	}
	if isBadReq {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	if req.Offset < 0 || req.Limit < 0 || req.Limit > 50 {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// search
	collections, err := db.SearchCollections(
		context.TODO(),
		logger,
		req.TimeRange,
		req.FloorPriceMin,
		req.FloorPriceMax,
		req.SaleCountMin,
		req.SaleCountMax,
		req.RoyaltyMin,
		req.RoyaltyMax,
		req.ProfitMarginMin,
		req.ProfitMarginMax,
		req.OwnerPercentageMin,
		req.OwnerPercentageMax,
		req.CollectionAgeMin,
		req.CollectionAgeMax,
		req.Keyword,
		req.Sort,
		req.Asc,
		req.Offset,
		req.Limit,
	)
	if err != nil {
		logger.Error("[DB] Failed to read cached collections")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	ctx.JSON(map[string]interface{}{"collections": collections})
}

func GetCollectionDetail(ctx iris.Context) {
	logger := ctx.Application().Logger()

	// parse params
	slug := ctx.Params().GetString("slug")
	if slug == "" {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// read index
	contract, _ := db.GetCollectionContract(context.TODO(), slug)
	if contract == "" {
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	// fetch info
	infoData, err := http.GetOpenSeaCollectionInfo(logger, slug)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch collection detail")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	var info map[string]interface{}
	err = json.Unmarshal([]byte(gjson.Get(infoData, "collection").String()), &info)
	if err != nil {
		logger.Error("[API] Failed to deserialize collection detail")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}
	// fetch collection offers
	offerData, err := http.GetOpenSeaCollectionOffers(logger, slug)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch collection offers")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	var offers []interface{}
	err = json.Unmarshal([]byte(gjson.Get(offerData, "offers").String()), &offers)
	if err != nil {
		logger.Error("[API] Failed to deserialize collection offers")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}
	// fetch listings
	listingData, err := http.GetReservoirCollectionListing(logger, contract)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch collection listings")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	// fetch recent sales
	salesData, err := http.GetOpenSeaCollectionRecentSales(logger, slug)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch collection sales")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	// fetch collection graph
	graphData, err := http.GetReservoirCollectionDaily(logger, contract)
	if err != nil {
		logger.Error("[HTTP] Failed to fetch collection graph")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	var graph []interface{}
	err = json.Unmarshal([]byte(gjson.Get(graphData, "collections").String()), &graph)
	if err != nil {
		logger.Error("[API] Failed to deserialize collection graph")
		ctx.StopWithStatus(iris.StatusInternalServerError)
	}

	// format data
	infoRes := map[string]interface{}{}
	offersRes := []map[string]interface{}{}
	listingsRes := []map[string]interface{}{}
	salesRes := []map[string]interface{}{}

	count := 0
	for _, offer := range gjson.Get(offerData, "offers").Array() {
		chain := gjson.Get(offer.String(), "chain").String()
		token := gjson.Get(offer.String(), "protocol_data.parameters.offer.0.token").String()
		if chain == "ethereum" && token == "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" {
			offersRes = append(offersRes, map[string]interface{}{
				"price":      gjson.Get(offer.String(), "protocol_data.parameters.offer.0.endAmount").Float() / 1000000000000000000,
				"start_time": gjson.Get(offer.String(), "protocol_data.parameters.startTime").Int(),
				"end_time":   gjson.Get(offer.String(), "protocol_data.parameters.endTime").Int(),
			})
			count++
		}
		if count >= 5 {
			break
		}
	}
	count = 0
	for _, listing := range gjson.Get(listingData, "orders").Array() {
		token := gjson.Get(listing.String(), "price.currency.symbol").String()
		if token == "ETH" {
			listingsRes = append(listingsRes, map[string]interface{}{
				"token_id":    gjson.Get(listing.String(), "criteria.data.token.tokenId").String(),
				"price":       gjson.Get(listing.String(), "price.amount.decimal").Float(),
				"valid_from":  gjson.Get(listing.String(), "validFrom").Int(),
				"valid_until": gjson.Get(listing.String(), "validUntil").Int(),
			})
			count++
		}
		if count >= 5 {
			break
		}
	}
	count = 0
	for _, sale := range gjson.Get(salesData, "asset_events").Array() {
		token := gjson.Get(sale.String(), "payment_token.address").String()
		if token == "0x0000000000000000000000000000000000000000" {
			salesRes = append(salesRes, map[string]interface{}{
				"token_id":  gjson.Get(sale.String(), "asset.token_id").String(),
				"image_url": gjson.Get(sale.String(), "asset.image_preview_url").String(),
				"price":     gjson.Get(sale.String(), "total_price").Float(),
				"date":      gjson.Get(sale.String(), "created_date").String(),
			})
			count++
		}
		if count >= 5 {
			break
		}
	}

	// update db
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
	if info["is_creator_fees_enforced"].(bool) {
		info["total_royalty"] = info["opensea_seller_fee_basis_points"].(int) + info["dev_seller_fee_basis_points"].(int)
	} else {
		info["total_royalty"] = info["opensea_seller_fee_basis_points"].(int)
	}
	if len(offersRes) > 0 {
		info["top_bid_price"] = offersRes[0]["price"]
	}
	info["last_updated"] = time.Now().Unix()
	err = db.UpdateCollectionDetail(context.TODO(), logger, info)
	if err != nil {
		logger.Error("[DB] Failed to update collection detail")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	// merge results
	infoRes["name"] = info["name"]
	infoRes["image_url"] = info["image_url"]
	infoRes["total_supply"] = info["stats"].(map[string]interface{})["total_supply"]
	infoRes["one_day_sales"] = info["stats"].(map[string]interface{})["one_day_sales"]
	infoRes["one_day_sales_change"] = info["stats"].(map[string]interface{})["one_day_sales_change"]
	infoRes["floor_price"] = info["stats"].(map[string]interface{})["floor_price"]
	infoRes["total_royalty"] = info["total_royalty"]
	infoRes["top_bid_price"] = info["top_bid_price"]

	infoRes["listings"] = listingsRes
	infoRes["collection_offers"] = offersRes
	infoRes["recent_sales"] = salesRes
	infoRes["graph"] = graph

	ctx.JSON(infoRes)
}
