package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	awsRegion string = "us-east-2"
	fbPostUrl string = "https://graph.facebook.com/%s/feed?access_token=%s"
)

var (
	runIntervalMinutes int64
	fbLLAT string
	client *http.Client
	svc *sns.SNS
	topicARN string
)

type FBGroupFeed struct {
	Data []FBGroupPost `json:"data"`
}

type FBGroupPost struct {
	Id string `json:"id"`
	Message string `json:"message"`
	UpdatedTime string `json:"updatedTime"`
}

type Post struct {
	Message string `json:"message"`
	PostTime string `json:"postTime"`
}

type RequestBody struct {
	GroupIds []string `json:"groupIds"`
}

// TODO: check if post already exists and return error
func handle(ctx context.Context, event events.CloudWatchEvent) (events.APIGatewayProxyResponse, error) {
	log.Println("Event body: ", event)
	log.Println("context ", ctx)
	headers := map[string]string{"Access-Control-Allow-Origin": "*", "Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept"}

	loadEnv();

	var body RequestBody
	jsonParseError := json.Unmarshal([]byte(event.Detail), &body)
	if jsonParseError != nil {
		log.Println(jsonParseError)
		return events.APIGatewayProxyResponse{500, headers, nil, "Internal Server Error", false}, nil
	}

	// TODO use AWS Lambda configuration params (esp. time)
	posts, err := getFbGroupPosts(body.GroupIds)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{500, headers, nil, "Internal Server Error", false}, nil
	}

	count, err := publishPostsToTopic(posts)

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

func publishPostsToTopic(posts []FBGroupPost) (int, error) {
	count := 0
	for _, post := range posts {
		err := publishPostToTopic(post)
		if err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func publishPostToTopic(post FBGroupPost) error {
	json, err := json.Marshal(post)

	if err != nil {
		return err
	}

	out, err := svc.Publish(&sns.PublishInput{
		Message:  aws.String(string(json)),
		TopicArn: aws.String(topicARN),
	})

	if err != nil {
		return err
	}

	log.Println("SNS response: ", out)

	return nil
}

func getFbGroupPosts(groupIds []string) ([]FBGroupPost, error) {
	var allPosts []FBGroupPost
	for groupId := range groupIds {
		groupPosts, err := getFbGroupPost(groupIds[groupId])
		if err != nil {
			return nil, err
		}
		allPosts = append(allPosts, groupPosts...)
	}
	return allPosts, nil
}

func getFbGroupPost(groupId string) ([]FBGroupPost, error) {
	req, err := http.NewRequest("GET", getFbPostUrl(groupId), nil);
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	log.Println("FB Response: ", resp)
	responseBody, err := ioutil.ReadAll(resp.Body)

	var feed FBGroupFeed
	err = json.Unmarshal([]byte(responseBody), &feed)
	if err != nil {
		return nil, err;
	}

	return feed.Data, nil
}

func getFbPostUrl(groupId string) string {
	return fmt.Sprintf(fbPostUrl, groupId, fbLLAT)
}

func loadEnv() {
	loadEnvVars();
	loadHttpClient();
}

func loadEnvVars() {
	fbLLAT = os.Getenv("FBLLAT")
	topicARN = os.Getenv("TopicARN")
	hours, err := strconv.ParseInt(os.Getenv("RunIntervalMinutes"), 10, 32)
	if err == nil {
		runIntervalMinutes = hours
	} else {
		log.Print(err)
	}
}

func loadHttpClient() {
	client = &http.Client{}
}

func main() {
	ses, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)
	svc = sns.New(ses)

	if err != nil {
		log.Println("Error initiating dynamodb for get_fb_group_posts lambda function ", err.Error())
	} else {
		log.Println("Successfully initiated dynamodb for get_fb_group_posts lambda function")
		lambda.Start(handle)
	}
}