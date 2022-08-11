package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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

const (
	ClassRangeUnit = "classes"
)

var (
	awsConfig aws.Config
	resources gocms.DynamoDBResources
)

type HandlerFunc func(context.Context, events.APIGatewayProxyRequest, *events.APIGatewayProxyResponse) (interface{}, error)

type Handlers struct {
	Repo gocms.Repository
}

func (h Handlers) GetHandler(request events.APIGatewayProxyRequest) (f HandlerFunc, found bool) {
	key := fmt.Sprintf("%s %s", request.HTTPMethod, request.Resource)
	funcMap := map[string]HandlerFunc{
		"GET /classes":  h.ClassList,
		"POST /classes": h.ClassCreate,
	}
	f, found = funcMap[key]
	return
}

func (h Handlers) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	response.StatusCode = http.StatusOK
	response.Headers = map[string]string{
		"Content-Type":                  "application/json",
		"Access-Control-Expose-Headers": "Content-Range, X-Total-Count",
	}

	var data interface{}
	if handler, found := h.GetHandler(request); found {
		data, err = handler(ctx, request, &response)
	} else {
		response.StatusCode = http.StatusNotFound
		err = errors.New("no handler found for resource: " + request.Resource)
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
	if os.Getenv("USER") == "localstack" {
		endpointResolverFunc := func(service string, region string, options ...interface{}) (endpoint aws.Endpoint, err error) {
			endpoint = aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://localhost:4566", // 4566 for LocalStack; 8000 for amazon/dynamodb-local
				SigningRegion: "us-east-1",             // Must be a legitimate region for LocalStack S3 to work
			}
			return
		}
		endpointResolver := aws.EndpointResolverWithOptionsFunc(endpointResolverFunc)

		awsConfig, err = config.LoadDefaultConfig(
			context.Background(),
			config.WithEndpointResolverWithOptions(endpointResolver),
		)
	} else {
		awsConfig, err = config.LoadDefaultConfig(context.Background())
	}
	if err != nil {
		log.Fatalf("Failed to load default config: %v", err)
	}

	resources.FromEnv()

	handlers := Handlers{
		Repo: gocms.NewDynamoDBRepository(awsConfig, resources),
	}
	lambda.Start(handlers.HandleRequest)
}

//
// Handlers
//

func (h Handlers) ClassCreate(ctx context.Context, request events.APIGatewayProxyRequest, response *events.APIGatewayProxyResponse) (value interface{}, err error) {
	var class gocms.Class
	reader := strings.NewReader(request.Body)
	if err = json.NewDecoder(reader).Decode(&class); err != nil {
		return
	}

	classService := gocms.NewClassService(h.Repo)
	if err = classService.Create(ctx, &class); err != nil {
		return
	}

	return class, nil
}

func (h Handlers) ClassList(ctx context.Context, request events.APIGatewayProxyRequest, response *events.APIGatewayProxyResponse) (value interface{}, err error) {
	classService := gocms.NewClassService(h.Repo)

	filter := gocms.ClassFilter{
		Range: gocms.Range{End: 9},
	}
	classes, r, err := classService.List(ctx, filter)
	if err != nil {
		return
	}

	response.Headers["Content-Range"] = r.ContentRangeHeader(ClassRangeUnit)
	response.Headers["X-Total-Count"] = fmt.Sprint(r.Size)
	return classes, nil
}
