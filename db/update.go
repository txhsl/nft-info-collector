package db

import (
	"context"

	"github.com/kataras/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	logger.Info("[DB] Trends deleted: ", delete.DeletedCount)

	// insert new
	insert, err := coll.InsertMany(ctx, collections)
	if err != nil {
		return err
	}
	logger.Info("[DB] Trends inserted: ", len(insert.InsertedIDs))
	return nil
}

func UpdateCollectionIndex(ctx context.Context, logger *golog.Logger, coll *mongo.Collection, collections []interface{}) error {
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
	logger.Info("[DB] Index matched: ", update.MatchedCount, ", upserted: ", update.UpsertedCount, ", modified: ", update.ModifiedCount, ", deleted: ", update.DeletedCount, ", inserted: ", update.InsertedCount)
	return nil
}

func UpdateCollectionDetails(ctx context.Context, logger *golog.Logger, coll *mongo.Collection, details []map[string]interface{}) error {
	models := []mongo.WriteModel{}
	for _, detail := range details {
		slug := detail["slug"]
		models = append(models, mongo.NewReplaceOneModel().SetUpsert(true).SetFilter(bson.M{"slug": slug}).SetReplacement(detail))
	}
	if len(models) == 0 {
		return nil
	}

	update, err := coll.BulkWrite(ctx, models)
	if err != nil {
		return err
	}
	logger.Info("[DB] Details matched: ", update.MatchedCount, ", upserted: ", update.UpsertedCount, ", modified: ", update.ModifiedCount, ", deleted: ", update.DeletedCount, ", inserted: ", update.InsertedCount)
	return nil
}

func UpdateCollectionMetrics(ctx context.Context, logger *golog.Logger, coll *mongo.Collection, collections []interface{}) error {
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
	logger.Info("[DB] Metrics matched: ", update.MatchedCount, ", upserted: ", update.UpsertedCount, ", modified: ", update.ModifiedCount, ", deleted: ", update.DeletedCount, ", inserted: ", update.InsertedCount)
	return nil
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
