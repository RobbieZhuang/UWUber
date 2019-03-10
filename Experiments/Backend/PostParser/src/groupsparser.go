package main

import (
	"encoding/json"
	"fmt"
) 




type Post struct {
	ID   string `json:"id"`
	Message string `json:"message"`
}



func processPosts() {
	var post Post
	sampleJsonString = {
		"id": "id1",
		"message": "Looking for ride to Union/Finch from Waterloo bk on Sunday (10th) after 4pm."
	}
	json.Unmarshal([byte(sampleJsonString), &post])
	fmt.Print(post.message)
}
