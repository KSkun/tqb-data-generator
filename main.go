package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
)

var client *mongo.Client

func initDatabase() {
	dbAddr, _ := os.LookupEnv("DB_ADDR")
	_client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(dbAddr))
	if err != nil {
		panic(err)
	}
	err = _client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		panic(err)
	}
	client = _client
}

func main() {
	//initDatabase()
	//initSceneMap()
	//initQuestionMap()
	initTopology()
	println(subjectMap)
}
