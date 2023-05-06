package db

import (
	"context"
	"math"
	"nft-info-collector/config"
	"time"

	"github.com/kataras/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, error) {
	options := options.Client().ApplyURI(config.Load().MongoDB.Url)
	return mongo.Connect(context.TODO(), options)
}

func getAscValue(asc bool) int {
	if asc {
		return 1
	}
	return -1
}

func getDetailTimePrefix(timeRange string) string {
	switch timeRange {
	case "1d":
		return "one_day_"
	case "7d":
		return "seven_day_"
	case "30d":
		return "thirty_day_"
	default:
		return "one_day_"
	}
}

func getMetricTimePrefix(timeRange string) string {
	switch timeRange {
	case "1d":
		return "1day"
	case "7d":
		return "7day"
	case "30d":
		return "30day"
	default:
		return "1d"
	}
}

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

	// match name & slug, last updated
	if flt != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"$text": bson.M{
					"$search": flt,
				},
			},
		})
	}
	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"stats.total_supply": bson.M{
				"$gt": 0,
			},
			"last_updated": bson.M{
				"$gte": time.Now().Add(-time.Hour * 24).Unix(),
			},
		},
	})

	// floor price
	if floorPriceMin != 0 || floorPriceMax != math.MaxInt {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"stats.floor_price": bson.M{
					"$gte": floorPriceMin,
					"$lte": floorPriceMax,
				},
			},
		})
	}

	// sale count
	if saleCountMin != 0 || saleCountMax != math.MaxInt {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"stats." + getDetailTimePrefix(timeRange) + "sales": bson.M{
					"$gte": saleCountMin,
					"$lte": saleCountMax,
				},
			},
		})
	}

	// collection age
	if collectionAgeMin != 0 || collectionAgemax != math.MaxInt {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"created_date": bson.M{
					"$gte": collectionAgeMin,
					"$lte": collectionAgemax,
				},
			},
		})
	}

	// join metrics
	pipeline = append(pipeline, bson.M{
		"$lookup": bson.M{
			"from":         "collection-metrics",
			"localField":   "slug",
			"foreignField": "slug",
			"as":           "metrics",
		},
	})
	pipeline = append(pipeline, bson.M{
		"$replaceRoot": bson.M{
			"newRoot": bson.M{
				"$mergeObjects": bson.A{"$$ROOT", bson.M{"$arrayElemAt": bson.A{"$metrics", 0}}},
			},
		},
	})
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"metrics": 0,
		},
	})

	// only erc721
	pipeline = append(pipeline, bson.M{
		"$match": bson.M{
			"contractKind": bson.M{
				"$eq": "erc721",
			},
		},
	})

	// add fields
	pipeline = append(pipeline, bson.M{
		"$addFields": bson.M{
			"total_royalty": bson.M{
				"$add": bson.A{"$opensea_seller_fee_basis_points", "$dev_seller_fee_basis_points"},
			},
			"owners_percentage": bson.M{
				"$multiply": bson.A{100, bson.M{
					"$divide": bson.A{"$stats.num_owners", "$stats.total_supply"},
				}},
			},
			"profit_margin": bson.M{
				"$subtract": bson.A{bson.M{
					"$multiply": bson.A{"$floorAsk.price.amount.decimal", bson.M{
						"$subtract": bson.A{1, bson.M{
							"$divide": bson.A{"$dev_seller_fee_basis_points", 10000},
						}},
					}},
				}, "$topBid.price.amount.decimal"},
			},
		},
	})

	// royalty
	if royaltyMin != 0 || royaltyMax != 10000 {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"total_royalty": bson.M{
					"$gte": royaltyMin,
					"$lte": royaltyMax,
				},
			},
		})
	}

	// owner percentage
	if ownerPercentageMin != 0 || ownerPercentageMax != 100 {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"owners_percentage": bson.M{
					"$gte": ownerPercentageMin,
					"$lte": ownerPercentageMax,
				},
			},
		})
	}

	// profit margin
	if profitMarginMin != math.MinInt || profitMarginMax != math.MaxInt {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"profit_margin": bson.M{
					"$gte": profitMarginMin,
					"$lte": profitMarginMax,
				},
			},
		})
	}

	// project
	pipeline = append(pipeline, bson.M{
		"$project": bson.M{
			"_id": 0,
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
