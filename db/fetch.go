package db

import (
	"context"

	"github.com/kataras/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetSortedCollectionIndex(ctx context.Context, logger *golog.Logger, keyword string, asc bool, offset int, limit int) ([]bson.M, error) {
	// connect db
	client, err := Connect()
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())
	coll := client.Database("nft-info-collector").Collection("collection-index")

	// get collections
	sort := bson.M{keyword: getAscValue(asc)}
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
	logger.Info("[DB] Index searched: ", len(results))
	return results, nil
}

func GetCollectionLastUpdated(ctx context.Context, coll *mongo.Collection, slug string) (int64, error) {
	filter := bson.M{"slug": slug}
	option := options.FindOne().SetProjection(bson.M{"last_updated": 1})
	result := coll.FindOne(ctx, filter, option)
	if result.Err() != nil {
		return 0, result.Err()
	}
	var collection bson.M
	if err := result.Decode(&collection); err != nil {
		return 0, err
	}
	return collection["last_updated"].(int64), nil
}

func GetOffersLastUpdated(ctx context.Context, coll *mongo.Collection, slug string) (int64, error) {
	filter := bson.M{"slug": slug}
	option := options.FindOne().SetProjection(bson.M{"last_updated": 1})
	result := coll.FindOne(ctx, filter, option)
	if result.Err() != nil {
		return 0, result.Err()
	}
	var collection bson.M
	if err := result.Decode(&collection); err != nil {
		return 0, err
	}
	return collection["last_updated"].(int64), nil
}

func GetCollectionContract(ctx context.Context, slug string) (string, error) {
	// connect db
	client, err := Connect()
	if err != nil {
		return "", err
	}
	defer client.Disconnect(context.Background())
	coll := client.Database("nft-info-collector").Collection("collection-index")

	// get collections
	filter := bson.M{
		"opensea_slug":  slug,
		"contract_type": "ERC721",
	}
	option := options.FindOne()
	result := coll.FindOne(ctx, filter, option)
	if result.Err() != nil {
		return "", result.Err()
	}
	var collection bson.M
	if err := result.Decode(&collection); err != nil {
		return "", err
	}
	return []interface{}(collection["contracts"].(primitive.A))[0].(string), nil
}
