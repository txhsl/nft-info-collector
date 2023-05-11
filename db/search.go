package db

import (
	"context"

	"github.com/kataras/golog"
	"go.mongodb.org/mongo-driver/bson"
)

func SearchCollections(
	ctx context.Context,
	logger *golog.Logger,
	timeRange string,
	floorPriceMin int,
	floorPriceMax int,
	saleCountMin int,
	saleCountMax int,
	royaltyMin int,
	royaltyMax int,
	profitMarginMin int,
	profitMarginMax int,
	ownerPercentageMin int,
	ownerPercentageMax int,
	collectionAgeMin int,
	collectionAgemax int,
	flt string,
	srt string,
	asc bool,
	offset int,
	limit int,
) ([]bson.M, error) {
	// connect db
	client, err := Connect()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())
	coll := client.Database("nft-info-collector").Collection("collection-details")

	// get collections
	pipeline := []bson.M{}

	// match name & slug
	if flt != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"$text": bson.M{
					"$search": flt,
				},
			},
		})
	}
	// total supply, schema name, collection offers enabled, floor price, sale count, collection age, royalty, profit margin, owner percentage
	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"stats.total_supply": bson.M{
				"$gt": 0,
			},
			"primary_asset_contracts.0.schema_name": bson.M{
				"$eq": "ERC721",
			},
			"is_collection_offers_enabled": true,
			"stats.floor_price": bson.M{
				"$gte": floorPriceMin,
				"$lte": floorPriceMax,
			},
			"stats." + getDetailTimePrefix(timeRange) + "sales": bson.M{
				"$gte": saleCountMin,
				"$lte": saleCountMax,
			},
			"created_date": bson.M{
				"$gte": collectionAgeMin,
				"$lte": collectionAgemax,
			},
			"total_royalty": bson.M{
				"$gte": royaltyMin,
				"$lte": royaltyMax,
			},
		},
	})

	// join metrics
	pipeline = append(pipeline, bson.M{
		"$lookup": bson.M{
			"from":         "collection-offers",
			"localField":   "slug",
			"foreignField": "slug",
			"as":           "collection_offers",
		},
	})
	pipeline = append(pipeline, bson.M{
		"$replaceRoot": bson.M{
			"newRoot": bson.M{
				"$mergeObjects": bson.A{"$$ROOT", bson.M{"$arrayElemAt": bson.A{"$collection_offers", 0}}},
			},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"collection_offers": 0,
		},
	})

	// pop top bid
	pipeline = append(pipeline, bson.M{
		"$addFields": bson.M{
			"top_offer": bson.M{
				"$first": "$offers",
			},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$addFields": bson.M{
			"top_bid": bson.M{
				"$first": "$top_offer.protocol_data.parameters.offer",
			},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$addFields": bson.M{
			"top_bid_price": bson.M{
				"$divide": bson.A{bson.M{
					"$toDecimal": "$top_bid.endAmount",
				}, 1000000000000000000},
			},
		},
	})

	// add fields
	pipeline = append(pipeline, bson.M{
		"$addFields": bson.M{
			"owners_percentage": bson.M{
				"$multiply": bson.A{100, bson.M{
					"$divide": bson.A{"$stats.num_owners", "$stats.total_supply"},
				}},
			},
			"profit_margin": bson.M{
				"$subtract": bson.A{
					bson.M{
						"$multiply": bson.A{"$stats.floor_price", bson.M{
							"$subtract": bson.A{1, bson.M{
								"$divide": bson.A{"$total_royalty", 10000},
							}},
						}},
					},
					"$top_bid_price",
				},
			},
		},
	})

	// owner percentage, profit margin
	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"owners_percentage": bson.M{
				"$gte": ownerPercentageMin,
				"$lte": ownerPercentageMax,
			},
			"profit_margin": bson.M{
				"$gte": profitMarginMin,
				"$lte": profitMarginMax,
			},
		},
	})

	// sort
	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"stats." + srt: getAscValue(asc),
		},
	})
	// skip
	pipeline = append(pipeline, bson.M{
		"$skip": offset,
	})
	// limit
	pipeline = append(pipeline, bson.M{
		"$limit": limit,
	})
	pipeline = append(pipeline, bson.M{
		"$replaceRoot": bson.M{
			"newRoot": bson.M{
				"$mergeObjects": bson.A{"$stats", "$$ROOT"},
			},
		},
	})
	// project
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id":                  0,
			"name":                 1,
			"slug":                 1,
			"image_url":            1,
			"total_supply":         1,
			"total_royalty":        1,
			"floor_price":          1,
			"total_volume":         1,
			"one_day_sales":        1,
			"one_day_sales_change": 1,
			"owners_percentage":    1,
			"top_bid_price":        1,
			"profit_margin":        1,
			"last_updated":         1,
		},
	})

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	// analysis result
	results := []bson.M{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	logger.Info("[DB] Index searched: ", len(results))
	return results, nil
}
