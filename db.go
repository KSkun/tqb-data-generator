package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
)

var client *mongo.Client

var (
	colQuestion *mongo.Collection
	colScene    *mongo.Collection
	colSubject  *mongo.Collection
)

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

	db := client.Database("tqb-backend")
	colQuestion = db.Collection("question")
	colScene = db.Collection("scene")
	colSubject = db.Collection("subject")
}

var (
	subjects = make([]*bson.M, 0)
	questions = make([]*Question, 0)
	scenes = make([]*Scene, 0)
)

func dfsGenerateData(node *Node, lastQuestion primitive.ObjectID) primitive.ObjectID {
	question := Question{}
	if len(node.QuestionLabel) > 0 {
		question = questionMap[node.QuestionLabel]
		question.ID = primitive.NewObjectID()
		for k, v := range node.NextNode {
			question.NextScene = append(question.NextScene, NextSceneObj{
				Scene:    dfsGenerateData(v, question.ID),
				Option:   k,
			})
		}
	}
	scene := sceneMap[node.SceneLabel]
	scene.ID = primitive.NewObjectID()
	scene.FromQuestion = lastQuestion
	scene.NextQuestion = question.ID

	questions = append(questions, &question)
	scenes = append(scenes, &scene)
	return scene.ID
}

func importSubject(subjectMap map[string]*Node) {
	for k, v := range subjectMap {
		subjectObj := bson.M{
			"abbr": k,
			"name": k,
			"start_scene": dfsGenerateData(v, primitive.NilObjectID),
		}
		subjects = append(subjects, &subjectObj)
	}

	for _, subject := range subjects {
		_, _ = colSubject.InsertOne(context.Background(), subject)
	}
	for _, question := range questions {
		_, _ = colQuestion.InsertOne(context.Background(), question)
	}
	for _, scene := range scenes {
		_, _ = colScene.InsertOne(context.Background(), scene)
	}
}
