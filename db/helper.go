package db

import (
	"context"
	"nft-info-collector/config"

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
