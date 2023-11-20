package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var db = dynamodb.New(session.Must(session.NewSession()))

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	shortenedUrl := req.QueryStringParameters["url"]

	originalUrl, err := getOriginalUrl(shortenedUrl)
	fmt.Printf("OriginalUrl: %s. Error: %s\n", originalUrl, err)
	if err != nil {
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers: map[string]string{
			"Location": originalUrl,
		},
		Body: "Redirecting to " + originalUrl,
	}
	fmt.Printf("Lambda Response: %s\n", response)
	return response, nil
}

func getOriginalUrl(shortenedUrl string) (string, error) {
	fmt.Println("Attempt to get item by shortenedUrl: " + shortenedUrl)
	item, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("urls"),
		Key: map[string]*dynamodb.AttributeValue{
			"shortenedUrl": {
				S: aws.String(shortenedUrl),
			},
		},
	})

	fmt.Printf("Fetched item: %s. Error: %s\n", item.String(), err)
	return *item.Item["originalUrl"].S, err
}
