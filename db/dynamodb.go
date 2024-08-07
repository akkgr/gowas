package db

import (
	"context"
	"errors"
	"gettenant/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DynamoDBClientInterface interface {
	GetTenant(ctx context.Context, event *models.Request) (*models.Response, error)
}

type DynamoDBClient struct {
	Client    dynamodbiface.DynamoDBAPI
	Tablename string
}

func (d *DynamoDBClient) GetTenant(ctx context.Context, event *models.Request) (*models.Response, error) {
	if event == nil {
		return nil, errors.New("received nil event")
	}

	if event.Id == "" {
		return nil, errors.New("id is required")
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String("Tenants"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(event.Id),
			},
		},
	}

	result, err := d.Client.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("Could not find '" + event.Id + "'")
	}

	tenant := new(models.Response)
	err = dynamodbattribute.UnmarshalMap(result.Item, tenant)
	if err != nil {
		return nil, errors.New("Failed to unmarshal Record, " + err.Error())
	}

	return tenant, nil
}
