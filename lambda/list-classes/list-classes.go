package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jbaikge/gocms"
)

const RangeUnit = "classes"

var (
	dynamoConfig aws.Config
	resources gocms.DynamoDBResources
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	response.Headers = map[string]string{
		"Content-Type":                  "application/json",
		"Access-Control-Expose-Headers": "Content-Range, X-Total-Count",
		"Access-Control-Allow-Origin":   "*",
	}
	response.StatusCode = http.StatusOK

	out, err := func() (out interface{}, err error) {
		filter := gocms.ClassFilter{
			Range: gocms.Range{
				End: 9,
			},
		}

		err = filter.Range.ParseParams(request.QueryStringParameters)
		if err != nil {
			return
		}

		repo := gocms.NewDynamoDBRepository(dynamoConfig, resources)
		service := gocms.NewClassService(repo)

		list, r, err := service.List(context.Background(), filter)
		if err == gocms.ErrBadRange {
			response.StatusCode = http.StatusRequestedRangeNotSatisfiable
		}
		if err != nil {
			return
		}

		response.Headers["Content-Range"] = r.ContentRangeHeader(RangeUnit)
		response.Headers["X-Total-Count"] = fmt.Sprint(r.Size)
		return list, nil
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
