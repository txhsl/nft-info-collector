package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kataras/iris/v12"
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
	loadConfig()

	app := iris.New()
	app.Get("/", hello)
	collectionAPI := app.Party("/collections")
	{
		collectionAPI.Get("/immediate", listImmediateCollections)
		collectionAPI.Get("/cached", listCachedCollections)
		collectionAPI.Get("/{address}", collectionInfo)
	}
	app.Listen(":8080")
}

func loadConfig() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}
}

func connectDB() *mongo.Client {
	options := options.Client().ApplyURI(config.MongoDB)
	client, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		panic(err)
	}
	if err = client.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}

	return client
}

func hello(ctx iris.Context) {
	ctx.WriteString("Info collector is working.")
}

// collections
func listImmediateCollections(ctx iris.Context) {
	// build request
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://restapi.nftscan.com/api/v2/statistics/ranking/trade?time=7d&sort_field=volume&sort_direction=desc&show_7d_trends=false", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-API-KEY", config.APIKey)

	// send request
	res, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// analysis response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	data := gjson.Get(string(body), "data").String()

	// cache result
	var collections []interface{}
	if err = json.Unmarshal([]byte(data), &collections); err != nil {
		panic(err)
	}
	dbClient := connectDB()
	coll := dbClient.Database("nft-info-collector").Collection("collections")
	cacheCollections(context.TODO(), coll, collections)

	ctx.WriteString(data)
}

func listCachedCollections(ctx iris.Context) {
	client := connectDB()
	coll := client.Database("nft-info-collector").Collection("collections")
	collections := getCachedCollections(context.TODO(), coll)

	result, err := json.Marshal(collections)
	if err != nil {
		panic(err)
	}
	ctx.WriteString(string(result))
}

func getCachedCollections(ctx context.Context, coll *mongo.Collection) []bson.M {
	// get collections
	sort := bson.D{{Key: "volume", Value: -1}}
	filter := bson.D{{}}
	option := options.Find()
	option.SetSort(sort)
	option.SetLimit(100)
	cursor, err := coll.Find(ctx, filter, option)
	if err != nil {
		panic(err)
	}
	// analysis result
	results := []bson.M{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	fmt.Println("Collections searched:", len(results))
	return results
}

func cacheCollections(ctx context.Context, coll *mongo.Collection, collections []interface{}) {
	// delete old
	filter := bson.D{{}}
	delete, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		panic(err)
	}
	fmt.Println("Collections deleted:", delete.DeletedCount)

	// insert new
	insert, err := coll.InsertMany(ctx, collections)
	if err != nil {
		panic(err)
	}
	fmt.Println("Collections inserted:", len(insert.InsertedIDs))
}

func collectionInfo(ctx iris.Context) {
	address := ctx.Params().Get("address")
	println("Get " + address)
}
