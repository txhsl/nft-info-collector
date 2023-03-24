package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey  string `yaml:"apikey"`
	MongoDB string `yaml:"mongodb"`
}

var config = Config{}

func main() {
	// collector config
	loadConfig()

	// iris config
	app := iris.New()
	app.UseRouter(recover.New())
	customLogger := logger.New(logger.Config{
		Status: true,
		Method: true,
		Path:   true,
	})
	app.Use(customLogger)

	// iris routes
	app.Get("/", hello)
	collectionAPI := app.Party("/collections")
	{
		collectionAPI.Get("/immediate", listImmediateCollections)
		collectionAPI.Get("/cached", listCachedCollections)
		collectionAPI.Get("/{address}", collectionInfo)
	}
	app.Listen(":8080")
}

// config
func loadConfig() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}
}

func hello(ctx iris.Context) {
	ctx.WriteString("Info collector is working.")
}

// handler
// collections
func listImmediateCollections(ctx iris.Context) {
	logger := ctx.Application().Logger()
	// build request
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://restapi.nftscan.com/api/v2/statistics/ranking/trade?time=7d&sort_field=volume&sort_direction=desc&show_7d_trends=false", nil)
	if err != nil {
		logger.Error("[API] Failed to build nftscan request")
		panic(err)
	}
	req.Header.Add("X-API-KEY", config.APIKey)

	// send request
	res, err := httpClient.Do(req)
	if err != nil {
		logger.Error("[API] Failed to send nftscan request")
		panic(err)
	}
	defer res.Body.Close()

	// analysis response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("[API] Failed to read nftscan response")
		panic(err)
	}
	data := gjson.Get(string(body), "data").String()

	// cache result
	var collections []interface{}
	if err = json.Unmarshal([]byte(data), &collections); err != nil {
		logger.Error("[API] Failed to serialize collections")
		panic(err)
	}
	dbClient, err := connectDB()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		panic(err)
	}
	defer dbClient.Disconnect(context.TODO())
	coll := dbClient.Database("nft-info-collector").Collection("collections")
	cacheCollections(context.TODO(), logger, coll, collections)

	ctx.WriteString(data)
}

func listCachedCollections(ctx iris.Context) {
	logger := ctx.Application().Logger()
	// search db
	client, err := connectDB()
	if err != nil {
		logger.Error("[DB] Failed to connect mongodb")
		panic(err)
	}
	defer client.Disconnect(context.TODO())
	coll := client.Database("nft-info-collector").Collection("collections")
	collections, err := getCachedCollections(context.TODO(), logger, coll)
	if err != nil {
		logger.Error("[API] Failed to read cached collections")
		panic(err)
	}
	// serialize collections
	result, err := json.Marshal(collections)
	if err != nil {
		logger.Error("[API] Failed to deserialize cached collections")
		panic(err)
	}
	ctx.WriteString(string(result))
}

func collectionInfo(ctx iris.Context) {
	logger := ctx.Application().Logger()
	address := ctx.Params().Get("address")
	logger.Warn("[API] Get " + address)
}

// database
func connectDB() (*mongo.Client, error) {
	options := options.Client().ApplyURI(config.MongoDB)
	return mongo.Connect(context.TODO(), options)
}

func getCachedCollections(ctx context.Context, logger *golog.Logger, coll *mongo.Collection) ([]bson.M, error) {
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

func cacheCollections(ctx context.Context, logger *golog.Logger, coll *mongo.Collection, collections []interface{}) error {
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
