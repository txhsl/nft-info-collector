package db

import (
	"context"
	"nft-info-collector/config"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func buildFilter(filter primitive.M) (primitive.M, error) {
	// convert string to bool
	if filter["has_rarity"] != nil {
		value, err := strconv.ParseBool(filter["has_rarity"].(string))
		if err != nil {
			return nil, err
		}
		filter["has_rarity"] = value
	}

	// ingore expired collections
	expiration := 24 * time.Hour
	updateLimit := (time.Now().Add(-expiration)).Unix()
	filter["last_updated"] = bson.M{"$gte": updateLimit}

	return filter, nil
}
