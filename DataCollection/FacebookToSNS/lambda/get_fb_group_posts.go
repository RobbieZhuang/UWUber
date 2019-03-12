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
	"time"
)

const tableName string = "job_list"
const awsRegion string = "us-east-1"

var svc *dynamodb.DynamoDB

type ResponseBody struct {
	Response string `json:"response"`
}

type RequestBody struct {
	UserId string `json:"userId"`
	Url string `json:"url"`
}

type JobEntity struct {
	UserId string `json:"userId"`
	Link string `json:"link"`
	Timestamp int `json:"timestamp"`
}

func saveToDatabase(job JobEntity) {
	dbMap, err := dynamodbattribute.MarshalMap(job)
	if err != nil {
		log.Println("Got error marshalling job entity:")
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
		log.Println("Successfully added ", job, " to table ", tableName)
		log.Println("Database response: ", out)
	}
}

// TODO: check if job already exists and return error
func handle(ctx context.Context, name events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Request body: ", name)
	log.Println("context ", ctx)
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept"}

	var body RequestBody
	jsonParseError := json.Unmarshal([]byte(name.Body), &body)
	if jsonParseError != nil {
		log.Println(jsonParseError)
		return events.APIGatewayProxyResponse{500, headers, "Internal Server Error", false}, nil
	}

	log.Println("Called by: ", body.UserId)
	code := 200
	response, jsonBuildError := json.Marshal(ResponseBody{Response: "Added " + body.Url + " to your Jobset!"})
	if jsonBuildError != nil {
		log.Println(jsonBuildError)
		response = []byte("Internal Server Error")
		code = 500
	}

	saveToDatabase(JobEntity{body.UserId, body.Url, int(time.Now().UnixNano() * int64(time.Nanosecond) / int64(time.Millisecond))})

	return events.APIGatewayProxyResponse{code, headers, string(response), false}, nil
}

func main() {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	// Create DynamoDB client
	svc = dynamodb.New(session)

	if err != nil {
		log.Println("Error initiating dynamodb for get_fb_group_posts lambda function ", err.Error())
	} else {
		log.Println("Successfully initiated dynamodb for get_fb_group_posts lambda function")
		lambda.Start(handle)
	}
}