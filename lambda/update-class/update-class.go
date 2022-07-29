package main

import (
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

var (
	dynamoConfig aws.Config
	resources gocms.DynamoDBResources
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	response.Headers = map[string]string{
		"Content-Type":                "application/json",
		"Access-Control-Allow-Origin": "*",
	}
	response.StatusCode = http.StatusOK

	out, err := func() (out interface{}, err error) {
		id, ok := request.PathParameters["id"]
		if !ok {
			response.StatusCode = http.StatusBadRequest
			err = fmt.Errorf("id not specified in URL")
			return
		}

		var class gocms.Class
		if err = json.NewDecoder(strings.NewReader(request.Body)).Decode(&class); err != nil {
			response.StatusCode = http.StatusBadRequest
			return
		}

		// Force ID to be what is in the URL. Not sure if necessary? Should prevent
		// changing a class ID.
		class.Id = id

		repo := gocms.NewDynamoDBRepository(dynamoConfig, resources)
		service := gocms.NewClassService(repo)

		if err = service.Update(ctx, &class); err != nil {
			response.StatusCode = http.StatusInternalServerError
			return
		}

		return class, nil
	}()

	if err != nil {
		out = struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}
	}

	body, err := json.Marshal(out)
	if err != nil {
		return
	}
	response.Body = string(body)

	return
}

func main() {
	var err error
	dynamoConfig, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load default config: %v", err)
	}

	resources.FromEnv()

	lambda.Start(HandleRequest)
}
