package db

import (
	"context"
	"nft-info-collector/config"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, error) {
	options := options.Client().ApplyURI(config.Load().MongoDB.Url)
	return mongo.Connect(context.TODO(), options)
}

func getSortValue(asc bool) int {
	if asc {
		return 1
	}
	return -1
}

func fitFilterType(filter primitive.M) (primitive.M, error) {
	if filter["has_rarity"] != nil {
		value, err := strconv.ParseBool(filter["has_rarity"].(string))
		if err != nil {
			return nil, err
		}
		filter["has_rarity"] = value
	}
	return filter, nil
}
