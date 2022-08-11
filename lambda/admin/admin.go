package main

import (
	"bytes"
	"context"
	"encoding/json"
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

type Handlers struct {
	Repo gocms.Repository
}

const (
	ClassRangeUnit = "classes"
)

var (
	awsConfig aws.Config
	resources gocms.DynamoDBResources
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	response.StatusCode = http.StatusOK
	response.Headers = map[string]string{
		"Content-Type":                  "application/json",
		"Access-Control-Expose-Headers": "Content-Range, X-Total-Count",
	}

	handlers := Handlers{
		Repo: gocms.NewDynamoDBRepository(awsConfig, resources),
	}

	var data interface{}
	switch fmt.Sprintf("%s %s", request.HTTPMethod, request.Resource) {
	case "GET /classes":
		data, err = handlers.ClassList(ctx, request, &response)
	case "POST /classes":
		data, err = handlers.ClassCreate(ctx, request, &response)
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

	lambda.Start(HandleRequest)
}

//
// Handlers
//

func (h Handlers) ClassCreate(ctx context.Context, request events.APIGatewayProxyRequest, response *events.APIGatewayProxyResponse) (class gocms.Class, err error) {
	reader := strings.NewReader(request.Body)
	if err = json.NewDecoder(reader).Decode(&class); err != nil {
		return
	}

	classService := gocms.NewClassService(h.Repo)
	err = classService.Create(ctx, &class)
	return
}

func (h Handlers) ClassList(ctx context.Context, request events.APIGatewayProxyRequest, response *events.APIGatewayProxyResponse) (classes []gocms.Class, err error) {
	classService := gocms.NewClassService(h.Repo)

	filter := gocms.ClassFilter{
		Range: gocms.Range{End: 9},
	}
	classes, r, err := classService.List(ctx, filter)

	response.Headers["Content-Range"] = r.ContentRangeHeader(ClassRangeUnit)
	response.Headers["X-Total-Count"] = fmt.Sprint(r.Size)
	return
}
