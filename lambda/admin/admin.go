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
	ClassRangeUnit    = "classes"
	DocumentRangeUnit = "documents"
)

var (
	awsConfig aws.Config
	resources gocms.DynamoDBResources
)

type HandlerFunc func(context.Context, events.APIGatewayV2HTTPRequest, *events.APIGatewayV2HTTPResponse) (interface{}, error)

type Handlers struct {
	Repo gocms.Repository
}

func (h Handlers) GetHandler(request events.APIGatewayV2HTTPRequest) (f HandlerFunc, found bool) {
	// key := fmt.Sprintf("%s %s", request.HTTPMethod, request.Resource)
	key := request.RouteKey
	funcMap := map[string]HandlerFunc{
		"GET /classes":                                  h.ClassList,
		"POST /classes":                                 h.ClassCreate,
		"GET /classes/{class_id}":                       h.ClassById,
		"PUT /classes/{class_id}":                       h.ClassUpdate,
		"DELETE /classes/{class_id}":                    h.ClassDelete,
		"GET /classes/{class_id}/documents":             h.DocumentList,
		"POST /classes/{class_id}/documents":            h.DocumentCreate,
		"GET /classes/{class_id}/documents/{doc_id}":    h.DocumentById,
		"PUT /classes/{class_id}/documents/{doc_id}":    h.DocumentUpdate,
		"DELETE /classes/{class_id}/documents/{doc_id}": h.DocumentDelete,
		"GET /documents/{doc_id}":                       h.DocumentById,
		"PUT /documents/{doc_id}":                       h.DocumentUpdate,
		"DELETE /documents/{doc_id}":                    h.DocumentDelete,
	}
	f, found = funcMap[key]
	return
}

func (h Handlers) HandleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	response.StatusCode = http.StatusOK
	response.Headers = map[string]string{
		"Content-Type":                  "application/json",
		"Access-Control-Expose-Headers": "Content-Range, X-Total-Count",
		"Access-Control-Allow-Origin":   "*",
	}

	var data interface{}
	if handler, found := h.GetHandler(request); found {
		data, err = handler(ctx, request, &response)
	} else {
		response.StatusCode = http.StatusNotFound
		err = errors.New("no handler found for resource: " + request.RouteKey)
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

func (h Handlers) ClassById(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	id, ok := request.PathParameters["class_id"]
	if !ok {
		response.StatusCode = http.StatusBadRequest
		return nil, errors.New("no class_id specified")
	}

	return gocms.NewClassService(h.Repo).ById(ctx, id)
}

func (h Handlers) ClassCreate(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	var class gocms.Class
	reader := strings.NewReader(request.Body)
	if err = json.NewDecoder(reader).Decode(&class); err != nil {
		return
	}

	if err = gocms.NewClassService(h.Repo).Create(ctx, &class); err != nil {
		return
	}

	return class, nil
}

func (h Handlers) ClassDelete(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	id, ok := request.PathParameters["class_id"]
	if !ok {
		response.StatusCode = http.StatusBadRequest
		return nil, errors.New("no class_id specified")
	}

	err = gocms.NewClassService(h.Repo).Delete(ctx, id)
	return
}

func (h Handlers) ClassList(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	filter := gocms.ClassFilter{
		Range: gocms.Range{End: 9},
	}
	classes, r, err := gocms.NewClassService(h.Repo).List(ctx, filter)
	if err != nil {
		return
	}

	response.Headers["Content-Range"] = r.ContentRangeHeader(ClassRangeUnit)
	response.Headers["X-Total-Count"] = fmt.Sprint(r.Size)
	return classes, nil
}

func (h Handlers) ClassUpdate(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	id, ok := request.PathParameters["class_id"]
	if !ok {
		response.StatusCode = http.StatusBadRequest
		return nil, errors.New("no class_id specified")
	}

	var class gocms.Class
	if err = json.NewDecoder(strings.NewReader(request.Body)).Decode(&class); err != nil {
		response.StatusCode = http.StatusBadRequest
		return nil, fmt.Errorf("bad json: %w", err)
	}

	// Force ID to be what is in the URL. Not sure if necessary? Should prevent
	// changing a class ID.
	class.Id = id
	if err = gocms.NewClassService(h.Repo).Update(ctx, &class); err != nil {
		response.StatusCode = http.StatusInternalServerError
		return nil, err
	}

	return class, nil
}

func (h Handlers) DocumentById(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	id, ok := request.PathParameters["doc_id"]
	if !ok {
		response.StatusCode = http.StatusBadRequest
		return nil, fmt.Errorf("no doc_id specified")
	}

	return gocms.NewDocumentService(h.Repo).ById(ctx, id)
}

func (h Handlers) DocumentCreate(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	var doc gocms.Document
	reader := strings.NewReader(request.Body)
	if err = json.NewDecoder(reader).Decode(&doc); err != nil {
		return
	}

	classId, hasClassId := request.PathParameters["class_id"]
	if !hasClassId && doc.ClassId == "" {
		return nil, fmt.Errorf("no class_id specified in URL or body")
	}

	// URL is the authority. Set/Override Class ID based on the URL if it exists
	if hasClassId {
		doc.ClassId = classId
	}

	if err = gocms.NewDocumentService(h.Repo).Create(ctx, &doc); err != nil {
		return
	}

	return doc, nil
}

func (h Handlers) DocumentDelete(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	id, ok := request.PathParameters["doc_id"]
	if !ok {
		response.StatusCode = http.StatusBadRequest
		return nil, fmt.Errorf("no doc_id specified")
	}

	err = gocms.NewDocumentService(h.Repo).Delete(ctx, id)
	return
}

func (h Handlers) DocumentList(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	filter := gocms.DocumentFilter{
		Range: gocms.Range{End: 9},
	}
	docs, r, err := gocms.NewDocumentService(h.Repo).List(ctx, filter)
	if err != nil {
		return
	}

	response.Headers["Content-Range"] = r.ContentRangeHeader(DocumentRangeUnit)
	response.Headers["X-Total-Count"] = fmt.Sprint(r.Size)
	return docs, nil
}

func (h Handlers) DocumentUpdate(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	id, ok := request.PathParameters["doc_id"]
	if !ok {
		response.StatusCode = http.StatusBadRequest
		return nil, fmt.Errorf("no doc_id specified")
	}

	var doc gocms.Document
	if err = json.NewDecoder(strings.NewReader(request.Body)).Decode(&doc); err != nil {
		response.StatusCode = http.StatusBadRequest
		return nil, fmt.Errorf("bad json: %w", err)
	}

	// Force ID to be what it is in the URL.
	doc.Id = id
	if err = gocms.NewDocumentService(h.Repo).Update(ctx, &doc); err != nil {
		response.StatusCode = http.StatusInternalServerError
		return nil, err
	}

	return doc, nil
}
