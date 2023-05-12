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
	if limit > 50 {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

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

	// update db
	// detail
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
	info["last_updated"] = time.Now().Unix()
	err = db.UpdateCollectionDetail(context.TODO(), logger, info)
	if err != nil {
		logger.Error("[DB] Failed to update collection detail")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
	// offers
	doc := map[string]interface{}{}
	doc["slug"] = slug
	if len(offers) > 5 {
		doc["offers"] = offers[:5]
	} else {
		doc["offers"] = offers
	}
	doc["last_updated"] = time.Now().Unix()
	err = db.UpdateCollectionOffer(context.TODO(), logger, doc)
	if err != nil {
		logger.Error("[DB] Failed to update collection offer")
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	// format data
	infoRes := map[string]interface{}{}
	offersRes := []map[string]interface{}{}
	listingsRes := []map[string]interface{}{}
	salesRes := []map[string]interface{}{}

	for i, offer := range gjson.Get(offerData, "offers").Array() {
		offersRes = append(offersRes, map[string]interface{}{
			"price": gjson.Get(offer.String(), "protocol_data.parameters.offer.0.endAmount").Float() / 1000000000000000000,
		})
		if i >= 4 {
			break
		}
	}
	for i, listing := range gjson.Get(listingData, "orders").Array() {
		listingsRes = append(listingsRes, map[string]interface{}{
			"price": gjson.Get(listing.String(), "price.amount.decimal").Float(),
		})
		if i >= 4 {
			break
		}
	}
	count := 0
	for _, sale := range gjson.Get(salesData, "asset_events").Array() {
		token := gjson.Get(sale.String(), "payment_token.address").String()
		if token == "0x0000000000000000000000000000000000000000" {
			salesRes = append(salesRes, map[string]interface{}{
				"price": gjson.Get(sale.String(), "total_price").Float(),
				"date":  gjson.Get(sale.String(), "created_date").String(),
			})
			count++
		}
		if count >= 5 {
			break
		}
	}

	// merge results
	infoRes["name"] = info["name"]
	infoRes["image_url"] = info["image_url"]
	infoRes["total_supply"] = info["stats"].(map[string]interface{})["total_supply"]
	infoRes["one_day_sales"] = info["stats"].(map[string]interface{})["one_day_sales"]
	infoRes["floor_price"] = info["stats"].(map[string]interface{})["floor_price"]
	infoRes["total_royalty"] = info["total_royalty"]
	infoRes["top_bid_price"] = offersRes[0]["price"]

	infoRes["listings"] = listingsRes
	infoRes["offers"] = offersRes
	infoRes["sales"] = salesRes
	infoRes["graph"] = graph

	ctx.JSON(infoRes)
}
