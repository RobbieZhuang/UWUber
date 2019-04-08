package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var svc *dynamodb.DynamoDB

const (
	tableName    string = "Trips"
	awsRegion    string = "us-east-2"
	functionName string = "add_trip"
)

func persistToDB(trip Trip) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	svc = dynamodb.New(session)
	dbMap, err := dynamodbattribute.MarshalMap(trip)
	if err != nil {
		log.Println("Error persisting the parsed trip to DB")
		log.Println(err.Error())
	}

	dbInput := &dynamodb.PutItemInput{
		Item:      dbMap,
		TableName: aws.String(tableName),
	}

	out, err := svc.PutItem(dbInput)

	if err != nil {
		log.Println("Error putting item ", err.Error())
	} else {
		log.Println("Database response: ", out)
	}
}

func handleEvent(ctx context.Context, event events.DynamoDBEvent) (events.APIGatewayProxyResponse, error) {
	for _, record := range event.Records {
		for key := range record.Change.Keys {
			result, err := svc.GetItem(&dynamodb.GetItemInput{
				TableName: aws.String(tableName),
				Key: map[string]*dynamodb.AttributeValue{
					"TripID": {
						N: aws.String(key),
					},
				},
			})
			parseJson(dynamodbattribute.Unmarshal(result.Item))
		}
	}
}
