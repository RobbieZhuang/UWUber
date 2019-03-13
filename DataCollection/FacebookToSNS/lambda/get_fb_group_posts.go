package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"os"
	"strconv"
)

const (
	awsRegion string = "us-east-2"
)

var (
	svc *dynamodb.DynamoDB
	RunIntervalMinutes int64
	FBLLAT string
)

type RequestBody struct {
	GroupIds []string `json:"groupIds"`
}

// TODO: check if post already exists and return error
func handle(ctx context.Context, event events.CloudWatchEvent) (events.APIGatewayProxyResponse, error) {
	log.Println("Event body: ", event)
	log.Println("context ", ctx)
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept"}

	loadEnvVars();

	var body RequestBody
	jsonParseError := json.Unmarshal([]byte(event.Detail), &body)
	if jsonParseError != nil {
		log.Println(jsonParseError)
		return events.APIGatewayProxyResponse{500, headers, nil, "Internal Server Error", false}, nil
	}

	// TODO use AWS Lambda configuration params (esp. time)
	count, err := getFBGroupPosts(body.GroupIds)

	var (
		response string
		code int
	)

	if err != nil {
		log.Println(err)
		response = "Internal Server Error"
		code = 500
	} else {
		response = fmt.Sprintf("Added %d posts to topic", count);
		code = 200
	}

	return events.APIGatewayProxyResponse{code, headers, nil,response, false}, nil
}

func getFBGroupPosts(groupIds []string) (int, error) {
	postsAdded := 0
	for groupId := range groupIds {
		count, err := getFBGroupPost(groupIds[groupId])
		if err != nil {
			return postsAdded, err
		}
		postsAdded += count
	}
	return postsAdded, nil
}

func getFBGroupPost(groupId string) (int, error) {
	postsAdded := 0
	return postsAdded, nil
}

func loadEnvVars() {
	FBLLAT = os.Getenv("FBLLAT")
	hours, err := strconv.ParseInt(os.Getenv("RunIntervalMinutes"), 10, 32)
	if err == nil {
		RunIntervalMinutes = hours
	} else {
		log.Print(err)
	}
}

func main() {
	ses, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	// Create DynamoDB client
	svc = dynamodb.New(ses)

	if err != nil {
		log.Println("Error initiating dynamodb for get_fb_group_posts lambda function ", err.Error())
	} else {
		log.Println("Successfully initiated dynamodb for get_fb_group_posts lambda function")
		lambda.Start(handle)
	}
}