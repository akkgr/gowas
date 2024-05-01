package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-xray-sdk-go/xray"
)

type Request struct {
	Id string `json:"id"`
}

type Response struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, event *Request) (*Response, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	if event.Id == "" {
		return nil, errors.New("id is required")
	}

	log.Printf("Processing request data for id %s.\n", event.Id)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)
	xray.AWS(svc.Client)

	result, err := svc.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("products"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(event.Id),
			},
		},
	})

	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		msg := "Could not find '" + event.Id + "'"
		return nil, errors.New(msg)
	}

	item := Response{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Fatalf("Failed to unmarshal Record, %v", err)
	}

	log.Printf("Found item: %v\n", item)
	return &item, nil
}

func main() {
	lambda.Start(HandleRequest)
}
