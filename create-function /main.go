package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
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
	originalUrl := req.QueryStringParameters["url"]
	shortenedUrl, err := createShortUrl(originalUrl)

	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Shortened URL: " + shortenedUrl,
	}, nil
}

func createShortUrl(originalUrl string) (string, error) {
	shortenedUrl := generateShortLink(originalUrl)

	item := map[string]*dynamodb.AttributeValue{
		"shortenedUrl": {
			S: aws.String(shortenedUrl),
		},
		"originalUrl": {
			S: aws.String(originalUrl),
		},
	}
	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("urls"),
	})

	return shortenedUrl, err
}

func generateShortLink(originalURL string) string {
	// Step 1: Hash the original URL
	hash := sha256.Sum256([]byte(originalURL))

	// Step 2: Convert the hash to a Base64 string for URL shortening
	shortHash := base64.URLEncoding.EncodeToString(hash[:])

	// Step 3: Truncate to a reasonable length
	shortLink := shortHash[:8]

	return shortLink
}
