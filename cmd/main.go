package main

import (
	"gettenant/db"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ddb := dynamodb.New(sess)
	xray.AWS(ddb.Client)

	dynamoDbClient := &db.DynamoDBClient{
		Client:    ddb,
		Tablename: "Tenants",
	}

	lambda.Start(dynamoDbClient.GetTenant)
}
