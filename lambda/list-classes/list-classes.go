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

const RangeKey = "classes"

var (
	dynamoConfig aws.Config
	dynamoTables gocms.DynamoDBTables
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	response.Headers = map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Headers": "Content-Range",
		"Access-Control-Allow-Origin":  "*",
	}
	response.StatusCode = http.StatusOK

	out, err := func() (out interface{}, err error) {
		filter := gocms.ClassFilter{
			Range: gocms.Range{
				Start: 0,
				End:   10,
			},
		}

		rangeHeader, ok := request.Headers["range"]
		if ok && strings.HasPrefix(rangeHeader, RangeKey) {
			_, err = fmt.Scanf(RangeKey+"=%d-%d", &filter.Range.Start, &filter.Range.End)
			if err != nil {
				err = fmt.Errorf("bad range header: %s - %v", rangeHeader, err)
				return
			}
		}

		repo := gocms.NewDynamoDBRepository(dynamoConfig, dynamoTables)
		service := gocms.NewClassService(repo)

		list, r, err := service.List(context.Background(), filter)
		if err != nil {
			return
		}

		response.Headers["Content-Range"] = r.ContentRangeHeader(RangeKey)
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

	dynamoTables.FromEnv()

	lambda.Start(HandleRequest)
}
