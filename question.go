package main

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"os"
)

type Question struct {
	ID          primitive.ObjectID `json:"-" bson:"_id"`
	Label       string             `json:"label" bson:"label"`
	Title       string             `json:"title" bson:"title"`
	Desc        string             `json:"desc" bson:"desc"`
	Statement   string             `json:"statement" bson:"statement"`
	SubQuestion bson.M             `json:"sub_question" bson:"sub_question"`
	Author      string             `json:"author" bson:"author"`
	Audio       string             `json:"audio" bson:"audio"`
	TimeLimit   int                `json:"time_limit" bson:"time_limit"`
	NextScene   []NextSceneObj     `json:"-" bson:"next_scene"`
}

type NextSceneObj struct {
	Scene    primitive.ObjectID `json:"-" bson:"scene"`
	SceneStr string             `json:"-" bson:"-"`
	Option   string             `json:"-" bson:"option"`
}

func (q *Question) NewObjectID() {
	q.ID = primitive.NewObjectID()
}

func (q *Question) UpdateObjectID() {
	for _, scene := range q.NextScene {
		scene.Scene, _ = primitive.ObjectIDFromHex(scene.SceneStr)
	}
}

var questionMap = make(map[string]Question, 0)

func loadAllQuestion(path string, fileList []os.FileInfo) {
	for _, fileInfo := range fileList {
		question := Question{}
		file, _ := os.Open(path + "\\" + fileInfo.Name())
		bytes, _ := ioutil.ReadAll(file)
		_ = json.Unmarshal(bytes, &question)
		questionMap[question.Label] = question
	}
}

func initQuestionMap() {
	questionMap = make(map[string]Question, 0)
	println("[question] working on question english")
	fileList, _ := ioutil.ReadDir("english")
	loadAllQuestion("english", fileList)
	println("[question] working on question math")
	fileList, _ = ioutil.ReadDir("math")
	loadAllQuestion("math", fileList)
	println("[question] working on question physics")
	fileList, _ = ioutil.ReadDir("physics")
	loadAllQuestion("physics", fileList)
}
