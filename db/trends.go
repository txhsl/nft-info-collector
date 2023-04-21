package db

import (
	"context"

	"github.com/kataras/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCachedTrends(ctx context.Context, logger *golog.Logger, offset int, limit int) ([]bson.M, error) {
	// connect db
	client, err := Connect()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())
	coll := client.Database("nft-info-collector").Collection("trends")

	// get collections
	sort := bson.D{{Key: "volume", Value: -1}}
	filter := bson.D{{}}
	option := options.Find()
	option.SetSort(sort)
	option.SetSkip(int64(offset))
	option.SetLimit(int64(limit))
	cursor, err := coll.Find(ctx, filter, option)
	if err != nil {
		return nil, err
	}

	// analysis result
	results := []bson.M{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	logger.Info("[DB] Collections searched:", len(results))
	return results, nil
}

func ReplaceCachedTrends(ctx context.Context, logger *golog.Logger, collections []interface{}) error {
	// connect db
	client, err := Connect()
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())
	coll := client.Database("nft-info-collector").Collection("trends")

	// delete old
	filter := bson.D{{}}
	delete, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	logger.Info("[DB] Collections deleted: ", delete.DeletedCount)

	// insert new
	insert, err := coll.InsertMany(ctx, collections)
	if err != nil {
		return err
	}
	logger.Info("[DB] Collections inserted: ", len(insert.InsertedIDs))
	return nil
}
