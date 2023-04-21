package db

import (
	"context"

	"github.com/kataras/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSortedCollections(ctx context.Context, logger *golog.Logger, timeRange string, keyword string, asc bool, offset int, limit int) ([]bson.M, error) {
	// connect db
	client, err := Connect()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())
	coll := client.Database("nft-info-collector").Collection("collections-" + timeRange)

	// get collections
	sort := bson.D{{Key: keyword, Value: getSortValue(asc)}}
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
	logger.Info("[DB] Collections searched: ", len(results))
	return results, nil
}

func GetFilteredCollections(ctx context.Context, logger *golog.Logger, timeRange string, flt string, value string, offset int, limit int) ([]bson.M, error) {
	// connect db
	client, err := Connect()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())
	coll := client.Database("nft-info-collector").Collection("collections-" + timeRange)

	// get collections
	sort := bson.D{{Key: "volume_eth", Value: -1}}
	filter, err := fitFilterType(bson.M{flt: value})
	if err != nil {
		return nil, err
	}
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
	logger.Info("[DB] Collections searched: ", len(results))
	return results, nil
}

func UpdateCachedCollections(ctx context.Context, logger *golog.Logger, coll *mongo.Collection, collections []interface{}) error {
	models := []mongo.WriteModel{}
	for _, collection := range collections {
		slug := collection.(map[string]interface{})["slug"]
		models = append(models, mongo.NewReplaceOneModel().SetUpsert(true).SetFilter(bson.M{"slug": slug}).SetReplacement(collection))
	}
	if len(models) == 0 {
		return nil
	}

	update, err := coll.BulkWrite(ctx, models)
	if err != nil {
		return err
	}
	logger.Info("[DB] Collections matched: ", update.MatchedCount, ", upserted: ", update.UpsertedCount, ", modified: ", update.ModifiedCount, ", deleted: ", update.DeletedCount, ", inserted: ", update.InsertedCount)
	return nil
}
