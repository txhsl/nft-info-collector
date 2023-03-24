package db

import (
	"context"
	"nft-info-collector/config"

	"github.com/kataras/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, error) {
	options := options.Client().ApplyURI(config.Load().MongoDB)
	return mongo.Connect(context.TODO(), options)
}

func GetCachedCollections(ctx context.Context, logger *golog.Logger, coll *mongo.Collection) ([]bson.M, error) {
	// get collections
	sort := bson.D{{Key: "volume", Value: -1}}
	filter := bson.D{{}}
	option := options.Find()
	option.SetSort(sort)
	option.SetLimit(100)
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

func CacheCollections(ctx context.Context, logger *golog.Logger, coll *mongo.Collection, collections []interface{}) error {
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
