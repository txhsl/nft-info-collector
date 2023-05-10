package db

import (
	"context"

	"github.com/kataras/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func UpdateCollectionOffers(ctx context.Context, logger *golog.Logger, coll *mongo.Collection, offers []map[string]interface{}) error {
	models := []mongo.WriteModel{}
	for _, collection := range offers {
		slug := collection["slug"]
		models = append(models, mongo.NewReplaceOneModel().SetUpsert(true).SetFilter(bson.M{"slug": slug}).SetReplacement(collection))
	}
	if len(models) == 0 {
		return nil
	}

	update, err := coll.BulkWrite(ctx, models)
	if err != nil {
		return err
	}
	logger.Info("[DB] Offers matched: ", update.MatchedCount, ", upserted: ", update.UpsertedCount, ", modified: ", update.ModifiedCount, ", deleted: ", update.DeletedCount, ", inserted: ", update.InsertedCount)
	return nil
}
