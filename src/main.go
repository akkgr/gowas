package main

import (
	"gettenant/handlers"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handlers.HandleRequest)
}
