package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jbaikge/gocms"
)

var (
	lambdaConfig aws.Config
	resources    gocms.DynamoDBResources
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(request)
	return
}

func main() {
	var err error
	lambdaConfig, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load default config: %v", err)
	}

	resources.FromEnv()

	lambda.Start(HandleRequest)
}
