package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

const (
	tableName string = "TripPosts"
	awsRegion string = "us-east-2"
	functionName string = "add_trip_post"
)

var svc *dynamodb.DynamoDB

type ResponseBody struct {
	Response string `json:"response"`
}

type Post struct {
	Id string `json:"id"`
	Message string `json:"message"`
	UpdatedTime string `json:"updatedTime"`
}

type PostEntity struct {
	PostId string `json:"PostId"`
	Message string `json:"Message"`
	PostTime string `json:"PostTime"`
}

func saveToDatabase(post PostEntity) {
	dbMap, err := dynamodbattribute.MarshalMap(post)
	if err != nil {
		log.Println("Got error marshalling post entity:")
		log.Println(err.Error())
	}

	dbInput := &dynamodb.PutItemInput{
		Item: dbMap,
		TableName: aws.String(tableName),
	}

	out, err := svc.PutItem(dbInput)

	if err != nil {
		log.Println("Error calling put item ", err.Error())
	} else {
		log.Println("Successfully added ", post, " to table ", tableName)
		log.Println("Database response: ", out)
	}
}

// TODO: check if post already exists and return error
func handle(ctx context.Context, snsEvent events.SNSEvent) (events.APIGatewayProxyResponse, error) {
	log.Println("Post: ", snsEvent)
	log.Println("context ", ctx)
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept"}

	var post Post
	jsonParseError := json.Unmarshal([]byte(snsEvent.Records[0].SNS.Message), &post)
	if jsonParseError != nil {
		log.Println(jsonParseError)
		return events.APIGatewayProxyResponse{500, headers, nil, "Internal Server Error", false}, nil
	}

	log.Println("Post received ", post)
	code := 200
	response, jsonBuildError := json.Marshal(ResponseBody{Response: "Added " + post.Message + " to database"})
	if jsonBuildError != nil {
		log.Println(jsonBuildError)
		response = []byte("Internal Server Error")
		code = 500
	}

	saveToDatabase(PostEntity{post.Id, post.Message, post.UpdatedTime})

	return events.APIGatewayProxyResponse{code, headers, nil, string(response), false}, nil
}

func main() {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	// Create DynamoDB client
	svc = dynamodb.New(session)

	if err != nil {
		log.Println("Error initiating dynamodb for " + functionName + " lambda function ", err.Error())
	} else {
		log.Println("Successfully initiated dynamodb for " + functionName + " lambda function")
		lambda.Start(handle)
	}
}