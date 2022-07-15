package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jbaikge/gocms"
)

var (
	dynamoConfig aws.Config
	dynamoTable  string
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	classId, ok := request.PathParameters["id"]
	if !ok {
		err = fmt.Errorf("id not specified in URL")
		return
	}

	repo := gocms.NewDynamoDBRepository(dynamoConfig, dynamoTable)
	service := gocms.NewClassService(repo)

	class, err := service.ById(context.Background(), classId)
	if err != nil {
		return
	}

	body, err := json.Marshal(&class)
	if err != nil {
		return
	}

	response.Headers = map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}
	response.StatusCode = http.StatusOK
	response.Body = string(body)
	return
}

func main() {
	var err error
	dynamoConfig, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load default config: %v", err)
	}

	if dynamoTable = os.Getenv("DYNAMODB_TABLE"); dynamoTable == "" {
		log.Fatalf("DYNAMODB_TABLE environment variable not set")
	}

	lambda.Start(HandleRequest)
}
