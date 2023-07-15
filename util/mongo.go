package util

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db = NewMongo()
	bg = context.Background()
)

func addr[T any](v T) *T {
	return &v
}

func NewMongo() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:StrongPassword@mongo.mridulganga.dev:27017/"))
	if err != nil {
		panic(err)
	}
	return client
}

func MUpsert(dbName, coll, id string, data map[string]any) bool {
	data["_id"] = id
	_, err := db.Database(dbName).Collection(coll).UpdateOne(bg, bson.M{"_id": id}, bson.M{"$set": data}, &options.UpdateOptions{
		Upsert: addr(true),
	})
	if err != nil {
		log.Errorf("errow while MUpsert %v", err)
	}
	return err == nil
}

func MGet(dbName, coll, id string) map[string]any {
	result := db.Database(dbName).Collection(coll).FindOne(bg, bson.M{"_id": id}, nil)
	output := map[string]any{}
	result.Decode(&output)
	if len(output) > 0 {
		return output
	}
	return nil
}

func MDelete(dbName, coll, id string) bool {
	_, err := db.Database(dbName).Collection(coll).DeleteOne(bg, bson.M{"_id": id}, nil)
	if err != nil {
		log.Errorf("errow while MDelete %v", err)
	}
	return err == nil
}

func MFind(dbName, coll string, query map[string]any, sort map[string]int) []map[string]any {

	findOpts := &options.FindOptions{}
	if sort != nil {
		findOpts = &options.FindOptions{
			Sort: sort,
		}
	}

	result, err := db.Database(dbName).Collection(coll).Find(bg, query, findOpts)
	if err != nil {
		log.Errorf("error while MList %v", err)
		return nil
	}
	outputs := []map[string]any{}
	result.All(bg, &outputs)
	return outputs
}
