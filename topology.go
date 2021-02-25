package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type Node struct {
	QuestionLabel string
	SceneLabel    string
	NextNode      map[string]*Node
}

var subjectMap = make(map[string]*Node, 0)
var toFind = regexp.MustCompile("[0-9.]+")

func dfsXMindNode(node map[string]interface{}) *Node {
	nowNode := Node{}
	titleStr := node["title"].(string)
	allMatch := toFind.FindAllString(titleStr, -1)
	nowNode.QuestionLabel = allMatch[0]
	nowNode.SceneLabel = allMatch[2]
	nowNode.NextNode = make(map[string]*Node, 0)
	for _, childNode := range node["children"].(map[string]interface{})["attached"].([]interface{}) {
		optionText := childNode.(map[string]interface{})["title"].(string)
		nextNode := childNode.(map[string]interface{})["children"].(map[string]interface{})["attached"].([]interface{})[0].(map[string]interface{})
		nowNode.NextNode[optionText] = dfsXMindNode(nextNode)
	}
	return &nowNode
}

func initTopology() {
	var jsonMap []map[string]interface{}
	jsonFile, _ := os.Open("content.json")
	jsonStr, _ := ioutil.ReadAll(jsonFile)
	_ = json.Unmarshal(jsonStr, &jsonMap)
	headerNode := jsonMap[0]["rootTopic"].(map[string]interface{})["children"].(map[string]interface{})["attached"].([]interface{})
	for _, subjectHeader := range headerNode {
		subjectName := subjectHeader.(map[string]interface{})["title"].(string)
		subjectName = subjectName[0:strings.Index(subjectName, "；")]
		headerNode := subjectHeader.(map[string]interface{})["children"].(map[string]interface{})["attached"].([]interface{})[0].(map[string]interface{})
		print("[topology] working on subject " + subjectName)
		subjectMap[subjectName] = dfsXMindNode(headerNode)
	}
}
