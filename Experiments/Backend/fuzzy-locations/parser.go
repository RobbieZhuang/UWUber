package main

import (
	"encoding/json"
	"regexp"
)

type Post struct {
	Id          string `json:"id"`
	Username    string `json:”username”`
	Message     string `json:"message"`
	UpdatedTime string `json:"updatedTime"`
}

func createShitPost(id string, username string, message string, updatedTime string) {
	//couldNotParse := "CAN NOT PARSE, HUMAN VERIFICATION REQUIRED"
	//We will do something so we can manually create these, populate db with raw or output to csv?
}

func createTrip(id string, username string, message string, updatedTime string) {

}

func createTripRequest(id string, username string, message string, updatedTime string) {
	//Eventually once we get more trips
}

func parseMessage(id string, username string, message string, updatedTime string) {
	pattern, _ := regexp.Compile("(DRIV|LOOK)")
	match := pattern.FindString(message)
	if match == "DRIV" {
		createTrip(id, username, message, updatedTime)
	} else if match == "LOOK" {
		createTripRequest(id, username, message, updatedTime)
	} else {
		createShitPost(id, message, username, updatedTime)
	}
}

func parseJson(post []byte) {
	var parsePost Post
	err := json.Unmarshal(post, &parsePost)
	parseMessage(parsePost.Id, parsePost.Username, parsePost.Message, parsePost.UpdatedTime)
}
