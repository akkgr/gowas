package db

import (
	"context"
	"errors"
	"gettenant/models"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
}

func (m *mockDynamoDBClient) GetItemWithContext(ctx context.Context, input *dynamodb.GetItemInput, option ...request.Option) (*dynamodb.GetItemOutput, error) {

	if *input.Key["id"].S == "0" {
		return nil, errors.New("Could not find '0'")
	}

	output := &dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String("1"),
			},
			"name": {
				S: aws.String("Test Tenant"),
			},
		},
	}
	return output, nil
}

func TestGetTenantOK(t *testing.T) {
	mockSvc := &mockDynamoDBClient{}
	dynamo := &DynamoDBClient{
		Client:    mockSvc,
		Tablename: "Tenants",
	}

	event := &models.Request{
		Id: "1",
	}

	tenant, err := dynamo.GetTenant(context.Background(), event)

	assert.Equal(t, &models.Response{Id: "1", Name: "Test Tenant"}, tenant)
	assert.Nil(t, err)
}

func TestGetTenantNoEvent(t *testing.T) {
	mockSvc := &mockDynamoDBClient{}
	dynamo := &DynamoDBClient{
		Client:    mockSvc,
		Tablename: "Tenants",
	}

	tenant, err := dynamo.GetTenant(context.Background(), nil)

	assert.Equal(t, errors.New("received nil event"), err)
	assert.Nil(t, tenant)
}

func TestGetTenantNoId(t *testing.T) {
	mockSvc := &mockDynamoDBClient{}
	dynamo := &DynamoDBClient{
		Client:    mockSvc,
		Tablename: "Tenants",
	}

	event := &models.Request{
		Id: "",
	}

	tenant, err := dynamo.GetTenant(context.Background(), event)

	assert.Equal(t, errors.New("id is required"), err)
	assert.Nil(t, tenant)
}

func TestGetTenantError(t *testing.T) {
	mockSvc := &mockDynamoDBClient{}
	dynamo := &DynamoDBClient{
		Client:    mockSvc,
		Tablename: "Tenants",
	}

	event := &models.Request{
		Id: "0",
	}

	tenant, err := dynamo.GetTenant(context.Background(), event)

	assert.Equal(t, errors.New("Could not find '0'"), err)
	assert.Nil(t, tenant)
}
