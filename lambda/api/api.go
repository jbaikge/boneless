package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jbaikge/boneless"
)

type Error struct {
	Error string `json:"error"`
}

type FilterParam struct {
	Ids    []string
	Fields map[string]string
}

func (f *FilterParam) UnmarshalJSON(data []byte) (err error) {
	fields := make(map[string]json.RawMessage)
	if err = json.Unmarshal(data, &fields); err != nil {
		return
	}

	if f.Fields == nil {
		f.Fields = make(map[string]string)
	}

	for key, value := range fields {
		switch key {
		case "id":
			err = json.Unmarshal(value, &f.Ids)
		default:
			var s string
			err = json.Unmarshal(value, &s)
			f.Fields[key] = s
		}
		if err != nil {
			return
		}
	}
	return nil
}

const (
	ClassRangeUnit    = "classes"
	DocumentRangeUnit = "documents"
	TemplateRangeUnit = "templates"
)

var (
	awsConfig aws.Config
	resources boneless.DynamoDBResources
)

type HandlerFunc func(context.Context, events.APIGatewayV2HTTPRequest, *events.APIGatewayV2HTTPResponse) (interface{}, error)

type Handlers struct {
	Repo boneless.Repository
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
		"POST /files":                                   h.FileCreate,
		"POST /files/url":                               h.FileUploadUrl,
		"GET /templates":                                h.TemplateList,
		"POST /templates":                               h.TemplateCreate,
		"GET /templates/{template_id}":                  h.TemplateById,
		"PUT /templates/{template_id}":                  h.TemplateUpdate,
		"DELETE /templates/{template_id}":               h.TemplateDelete,
	}
	f, found = funcMap[key]
	return
}

func (h Handlers) HandleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	start := time.Now()
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
	response.Headers["X-Handler-Time"] = time.Since(start).String()
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
		Repo: boneless.NewDynamoDBRepository(awsConfig, resources),
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

	return boneless.NewClassService(h.Repo).ById(ctx, id)
}

func (h Handlers) ClassCreate(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	var class boneless.Class
	reader := strings.NewReader(request.Body)
	if err = json.NewDecoder(reader).Decode(&class); err != nil {
		return
	}

	if err = boneless.NewClassService(h.Repo).Create(ctx, &class); err != nil {
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

	err = boneless.NewClassService(h.Repo).Delete(ctx, id)
	return
}

func (h Handlers) ClassList(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	filter := boneless.ClassFilter{
		Range: boneless.Range{End: 9},
	}
	classes, r, err := boneless.NewClassService(h.Repo).List(ctx, filter)
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

	var class boneless.Class
	if err = json.NewDecoder(strings.NewReader(request.Body)).Decode(&class); err != nil {
		response.StatusCode = http.StatusBadRequest
		return nil, fmt.Errorf("bad json: %w", err)
	}

	// Force ID to be what is in the URL. Not sure if necessary? Should prevent
	// changing a class ID.
	class.Id = id
	if err = boneless.NewClassService(h.Repo).Update(ctx, &class); err != nil {
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

	return boneless.NewDocumentService(h.Repo).ById(ctx, id)
}

func (h Handlers) DocumentCreate(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	var doc boneless.Document
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

	if err = boneless.NewDocumentService(h.Repo).Create(ctx, &doc); err != nil {
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

	err = boneless.NewDocumentService(h.Repo).Delete(ctx, id)
	return
}

// filter: {} - For filtering, {"field":"value"}; for getMany, {"id":[1,2,3]}
// range: [0,9]
// sort: ["id","ASC"]
func (h Handlers) DocumentList(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	documentService := boneless.NewDocumentService(h.Repo)

	filter := boneless.DocumentFilter{
		Range: boneless.Range{End: 9},
	}

	if classId, ok := request.PathParameters["class_id"]; ok {
		filter.ClassId = classId
	}

	if param, ok := request.QueryStringParameters["range"]; ok {
		values := make([]int, 0, 2)
		if err = json.Unmarshal([]byte(param), &values); err != nil {
			return nil, fmt.Errorf("unmarshalling range %s: %w", param, err)
		}
		if len(values) != 2 {
			return nil, fmt.Errorf("not sure what to do with this range: %s", param)
		}
		filter.Range.Start = values[0]
		filter.Range.End = values[1]
	}

	if param, ok := request.QueryStringParameters["sort"]; ok {
		values := make([]string, 0, 2)
		if err = json.Unmarshal([]byte(param), &values); err != nil {
			return nil, fmt.Errorf("unmarshalling sort %s: %w", param, err)
		}
		if len(values) != 2 {
			return nil, fmt.Errorf("not sure what to do with this sort: %s", param)
		}
		filter.Sort.Field = strings.Replace(values[0], "values.", "", 1)
		filter.Sort.Direction = values[1]
	}

	// simple rest data provider calls "getMany" by using ?filter={"id":[1, 2, 3]}
	filterParam := new(FilterParam)
	if param, ok := request.QueryStringParameters["filter"]; ok {
		if err = json.Unmarshal([]byte(param), filterParam); err != nil {
			return nil, fmt.Errorf("unmarshalling filter parameter: %w", err)
		}
	}

	if len(filterParam.Ids) > 0 {
		docs := make([]boneless.Document, 0, len(filterParam.Ids))
		for _, id := range filterParam.Ids {
			doc, err := documentService.ById(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("getting documents by id: %w", err)
			}
			docs = append(docs, doc)
		}
		return docs, nil
	}

	// Handle remaining GET calls

	docs, r, err := documentService.List(ctx, filter)
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

	var doc boneless.Document
	if err = json.NewDecoder(strings.NewReader(request.Body)).Decode(&doc); err != nil {
		response.StatusCode = http.StatusBadRequest
		return nil, fmt.Errorf("bad json: %w", err)
	}

	// Force ID to be what it is in the URL.
	doc.Id = id
	if err = boneless.NewDocumentService(h.Repo).Update(ctx, &doc); err != nil {
		response.StatusCode = http.StatusInternalServerError
		return nil, err
	}

	return doc, nil
}

// This should handle file uploads for both documents and TinyMCE
// The latter expects a JSON document like the following:
// { "location": "folder/sub-folder/new-location.png" }
func (h Handlers) FileCreate(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	var reader io.Reader
	reader = strings.NewReader(request.Body)
	if request.IsBase64Encoded {
		reader = base64.NewDecoder(base64.StdEncoding, reader)
	}

	r, err := http.NewRequest(request.RequestContext.HTTP.Method, request.RequestContext.HTTP.Path, reader)
	if err != nil {
		err = fmt.Errorf("creating request: %w", err)
		return
	}
	for k, v := range request.Headers {
		r.Header.Set(k, v)
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		err = fmt.Errorf("getting form file: %w", err)
		return
	}
	defer file.Close()

	uploadFile := &boneless.File{
		Filename: fileHeader.Filename,
		Data:     file,
	}

	contentTypeBuffer := make([]byte, 512)
	if _, err = file.Read(contentTypeBuffer); err != nil {
		err = fmt.Errorf("reading first 512 bytes: %w", err)
		return
	}
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		err = fmt.Errorf("seeking to beginning of file: %w", err)
	}

	uploadFile.ContentType = http.DetectContentType(contentTypeBuffer)

	if err = h.Repo.CreateFile(ctx, uploadFile); err != nil {
		err = fmt.Errorf("creating file: %w", err)
		return
	}

	data := struct {
		Location string `json:"location"`
	}{
		Location: uploadFile.Location,
	}

	return data, nil
}

// Returns a Signed S3 URL with PUT access for the requestor to then PUT data to
func (h Handlers) FileUploadUrl(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	var uploadRequest boneless.FileUploadRequest
	if err = json.NewDecoder(strings.NewReader(request.Body)).Decode(&uploadRequest); err != nil {
		response.StatusCode = http.StatusBadRequest
		return nil, fmt.Errorf("bad json: %w", err)
	}

	return boneless.NewFileService(h.Repo).UploadUrl(ctx, uploadRequest)
}

func (h Handlers) TemplateById(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	id, ok := request.PathParameters["template_id"]
	if !ok {
		response.StatusCode = http.StatusBadRequest
		return nil, fmt.Errorf("no template_id specified")
	}

	return boneless.NewTemplateService(h.Repo).ById(ctx, id)
}

func (h Handlers) TemplateCreate(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	template := new(boneless.Template)
	reader := strings.NewReader(request.Body)
	if err = json.NewDecoder(reader).Decode(template); err != nil {
		return
	}

	if err = boneless.NewTemplateService(h.Repo).Create(ctx, template); err != nil {
		return
	}

	return template, nil
}

func (h Handlers) TemplateDelete(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	id, ok := request.PathParameters["template_id"]
	if !ok {
		response.StatusCode = http.StatusBadRequest
		return nil, fmt.Errorf("no template_id specified")
	}

	err = boneless.NewTemplateService(h.Repo).Delete(ctx, id)
	return
}

func (h Handlers) TemplateList(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	templateService := boneless.NewTemplateService(h.Repo)

	filter := boneless.TemplateFilter{
		Range: boneless.Range{End: 9},
	}

	if param, ok := request.QueryStringParameters["range"]; ok {
		values := make([]int, 0, 2)
		if err = json.Unmarshal([]byte(param), &values); err != nil {
			return nil, fmt.Errorf("unmarshalling range %s: %w", param, err)
		}
		if len(values) != 2 {
			return nil, fmt.Errorf("not sure what to do with this range: %s", param)
		}
		filter.Range.Start = values[0]
		filter.Range.End = values[1]
	}

	// simple rest data provider calls "getMany" by using ?filter={"id":[1, 2, 3]}
	filterParam := new(FilterParam)
	if param, ok := request.QueryStringParameters["filter"]; ok {
		if err = json.Unmarshal([]byte(param), filterParam); err != nil {
			return nil, fmt.Errorf("unmarshalling filter parameter: %w", err)
		}
	}

	if len(filterParam.Ids) > 0 {
		templates := make([]boneless.Template, 0, len(filterParam.Ids))
		for _, id := range filterParam.Ids {
			template, err := templateService.ById(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("getting templates by id: %w", err)
			}
			templates = append(templates, template)
		}
		return templates, nil
	}

	templates, r, err := templateService.List(ctx, filter)
	if err != nil {
		return
	}

	response.Headers["Content-Range"] = r.ContentRangeHeader(TemplateRangeUnit)
	response.Headers["X-Total-Count"] = fmt.Sprint(r.Size)
	return templates, nil
}

func (h Handlers) TemplateUpdate(ctx context.Context, request events.APIGatewayV2HTTPRequest, response *events.APIGatewayV2HTTPResponse) (value interface{}, err error) {
	id, ok := request.PathParameters["template_id"]
	if !ok {
		response.StatusCode = http.StatusBadRequest
		return nil, fmt.Errorf("no template_id specified")
	}

	template := new(boneless.Template)
	if err = json.NewDecoder(strings.NewReader(request.Body)).Decode(template); err != nil {
		response.StatusCode = http.StatusBadRequest
		return nil, fmt.Errorf("bad json: %w", err)
	}

	// Force ID to be what it is in the URL
	template.Id = id
	if err = boneless.NewTemplateService(h.Repo).Update(ctx, template); err != nil {
		response.StatusCode = http.StatusInternalServerError
		return nil, fmt.Errorf("update error: %w", err)
	}

	return template, nil
}
