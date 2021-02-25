package main

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"os"
)

type Scene struct {
	ID              primitive.ObjectID `json:"-" bson:"_id"`
	Label           string             `json:"label" bson:"label"`
	Title           string             `json:"title" bson:"title"`
	Text            string             `json:"text" bson:"text"`
	BGM             string             `json:"bgm" bson:"bgm"`
	FromQuestion    primitive.ObjectID `json:"-" bson:"from_question"`
	FromQuestionStr string             `json:"-" bson:"-"`
	NextQuestion    primitive.ObjectID `json:"-" bson:"next_question"`
	NextQuestionStr string             `json:"-" bson:"-"`
}

func (s *Scene) NewObjectID() {
	s.ID = primitive.NewObjectID()
}

func (s *Scene) UpdateObjectID() {
	s.FromQuestion, _ = primitive.ObjectIDFromHex(s.FromQuestionStr)
	s.NextQuestion, _ = primitive.ObjectIDFromHex(s.NextQuestionStr)
}

var sceneMap = make(map[string]Scene, 0)

func loadAllScene(path string, fileList []os.FileInfo) {
	for _, fileInfo := range fileList {
		scene := Scene{}
		file, _ := os.Open(path + "\\" + fileInfo.Name())
		bytes, _ := ioutil.ReadAll(file)
		_ = json.Unmarshal(bytes, &scene)
		sceneMap[scene.Label] = scene
	}
}

func initSceneMap() {
	sceneMap = make(map[string]Scene, 0)
	println("[scene] working on scene 英语 1")
	fileList, _ := ioutil.ReadDir("scene\\英语\\1")
	loadAllScene("scene\\英语\\1", fileList)
	println("[scene] working on scene 英语 2")
	fileList, _ = ioutil.ReadDir("scene\\英语\\2")
	loadAllScene("scene\\英语\\2", fileList)
}
