package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jbaikge/gocms"
)

const RangeUnit = "classes"

var (
	dynamoConfig aws.Config
	dynamoTables gocms.DynamoDBTables
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	response.Headers = map[string]string{
		"Content-Type":                  "application/json",
		"Access-Control-Expose-Headers": "Content-Range",
		"Access-Control-Allow-Origin":   "*",
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
		if ok && strings.HasPrefix(rangeHeader, RangeUnit) {
			if err = filter.Range.ParseHeader(rangeHeader, RangeUnit); err != nil {
				response.StatusCode = http.StatusRequestedRangeNotSatisfiable
				return
			}
		}

		rangeParam, ok := request.QueryStringParameters["range"]
		if filter.Range.IsZero() && ok {
			bounds := make([]int, 0, 2)
			err = json.Unmarshal([]byte(rangeParam), &bounds)
			if err != nil || len(bounds) != 2 {
				response.StatusCode = http.StatusRequestedRangeNotSatisfiable
				return
			}
			filter.Range.Start = bounds[0]
			filter.Range.End = bounds[1]
		}

		repo := gocms.NewDynamoDBRepository(dynamoConfig, dynamoTables)
		service := gocms.NewClassService(repo)

		list, r, err := service.List(context.Background(), filter)
		if err == gocms.ErrBadRange {
			response.StatusCode = http.StatusRequestedRangeNotSatisfiable
		}
		if err != nil {
			return
		}

		response.Headers["Content-Range"] = r.ContentRangeHeader(RangeUnit)
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
