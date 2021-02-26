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
	NextNodeText  string
	NextNode      map[string]*Node
}

var subjectMap = make(map[string]*Node, 0)
var toFind = regexp.MustCompile("[0-9.]+")

func nextXMindNode(node map[string]interface{}) []interface{} {
	if _, found := node["children"]; !found {
		return make([]interface{}, 0)
	}
	return node["children"].(map[string]interface{})["attached"].([]interface{})
}

func getTitleFromNode(node map[string]interface{}) string {
	return node["title"].(string)
}

func dfsXMindNode(node map[string]interface{}) *Node {
	nowNode := Node{}
	titleStr := getTitleFromNode(node)
	allMatch := toFind.FindAllString(titleStr, -1)
	nowNode.QuestionLabel = allMatch[0]
	nowNode.SceneLabel = allMatch[2]
	nowNode.NextNode = make(map[string]*Node, 0)

	childNodeTitle := getTitleFromNode(nextXMindNode(node)[0].(map[string]interface{}))
	if strings.Contains(childNodeTitle, "End") { // 到达 Ending
		sceneID := toFind.FindString(childNodeTitle)
		sceneText := sceneMap[sceneID].Text
		nextNode := nextXMindNode(node)[0] // Ending 节点
		flagHasAS := false
		if len(nextXMindNode(nextNode.(map[string]interface{}))) > 0 { // 有 After Story
			nextNextNode := nextXMindNode(nextNode.(map[string]interface{}))[0] // After Story 节点
			nextNextNextNode := nextXMindNode(nextNextNode.(map[string]interface{})) // AS 后续节点
			if len(nextNextNextNode) > 0 {
				flagHasAS = true
				sceneText += "\\n\\n*" + getTitleFromNode(nextNode.(map[string]interface{})) + "*"
				nowNode.NextNodeText = sceneText
				node = nextNextNode.(map[string]interface{})
			}
		}
		if !flagHasAS { // 无 AS
			// 建立仅包含 Ending 剧情的节点，终止递归
			nextNode := Node{}
			nextNode.SceneLabel = sceneID
			nowNode.NextNode["Ending"] = &nextNode
			return &nowNode
		}
	} else if strings.Contains(childNodeTitle, "；") { // 选择页剧情
		sceneID := toFind.FindString(childNodeTitle)
		sceneText := sceneMap[sceneID].Text + "\\n\\n*" + childNodeTitle[0:strings.Index(childNodeTitle, "；")] + "*"
		nowNode.NextNodeText = sceneText
		node = nextXMindNode(node)[0].(map[string]interface{})
	}

	for _, childNodeObj := range nextXMindNode(node) {
		childNode := childNodeObj.(map[string]interface{})
		optionText := getTitleFromNode(childNode)
		nextNode := nextXMindNode(childNode)[0].(map[string]interface{})
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
