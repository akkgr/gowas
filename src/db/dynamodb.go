package db

import (
	"context"
	"errors"
	"gettenant/models"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func GetTenant(ctx context.Context, event *models.Request) (*models.Response, error) {
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

	item := models.Response{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Fatalf("Failed to unmarshal Record, %v", err)
	}

	return &item, nil
}
