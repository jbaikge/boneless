package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jbaikge/gocms"
)

type Error struct {
	Error string `json:"error"`
}

var (
	awsConfig aws.Config
	resources gocms.DynamoDBResources
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	response.StatusCode = http.StatusOK
	response.Headers = map[string]string{
		"Content-Type": "application/json",
	}

	repo := gocms.NewDynamoDBRepository(awsConfig, resources)

	var data interface{}
	switch fmt.Sprintf("%s %s", request.HTTPMethod, request.Resource) {
	case "GET /classes":
		data, err = handleClassList(ctx, request, repo)
	case "POST /classes":
		data, err = handleClassCreate(ctx, request, repo)
	}

	// Redirect the error so it comes out as JSON instead of a 500
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		data = Error{Error: err.Error()}
	}

	var buffer bytes.Buffer
	if err = json.NewEncoder(&buffer).Encode(data); err != nil {
		return
	}
	response.Body = buffer.String()
	return
}

func main() {
	var err error
	awsConfig, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load default config: %v", err)
	}

	resources.FromEnv()

	lambda.Start(HandleRequest)
}

//
// Handlers
//

func handleClassCreate(ctx context.Context, request events.APIGatewayProxyRequest, repo gocms.Repository) (class gocms.Class, err error) {
	reader := strings.NewReader(request.Body)
	if err = json.NewDecoder(reader).Decode(&class); err != nil {
		return
	}

	classService := gocms.NewClassService(repo)
	err = classService.Create(ctx, &class)
	return
}

func handleClassList(ctx context.Context, request events.APIGatewayProxyRequest, repo gocms.Repository) (classes []gocms.Class, err error) {
	return
}
