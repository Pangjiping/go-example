package mongo_go

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func getCollection() (*mongo.Collection, error) {
	client, err := mongo.Connect(context.TODO(),
		options.Client().ApplyURI("mongodb://xxx.xxx.xxx.xxx:27017"),
		options.Client().SetConnectTimeout(time.Second*5))
	if err != nil {
		return nil, err
	}

	// 选择数据库
	database := client.Database("my_db")

	// 选择表
	collection := database.Collection("my_collection")
	return collection, nil
}
