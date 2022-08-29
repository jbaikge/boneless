package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jbaikge/gocms"
)

var (
	awsConfig aws.Config
	resources gocms.DynamoDBResources
)

type Frontend struct {
	Repo gocms.Repository
}

func (frontend Frontend) HandleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	documentService := gocms.NewDocumentService(frontend.Repo)
	document, err := documentService.ByPath(ctx, request.RawPath)
	if err != nil {
		response.StatusCode = http.StatusNotFound
		response.Body = "Not found!"
		return
	}

	templateService := gocms.NewTemplateService(frontend.Repo)
	filter := gocms.TemplateFilter{
		Range: gocms.Range{
			End: 1000,
		},
	}
	templates, _, err := templateService.List(ctx, filter)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = fmt.Sprintf("Something went wrong trying to fetch templates: %v", err)
		return
	}

	// Need to convert the target template ID to a name string
	var name string
	t := template.New(document.Id)
	for _, tmpl := range templates {
		if document.TemplateId == tmpl.Id {
			name = tmpl.Name
		}
		if _, err = t.New(tmpl.Name).Parse(tmpl.Body); err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Body = fmt.Sprintf("Template compilation error: %v", err)
			return
		}
	}

	data := struct {
		Document gocms.Document
	}{
		Document: document,
	}

	buffer := new(bytes.Buffer)
	t.ExecuteTemplate(buffer, name, data)
	response.StatusCode = http.StatusOK
	response.Headers = map[string]string{
		"Content-Type": "text/html",
	}
	response.Body = buffer.String()
	return
}

func main() {
	var err error
	awsConfig, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load default config: %v", err)
	}

	resources.FromEnv()

	frontend := Frontend{
		Repo: gocms.NewDynamoDBRepository(awsConfig, resources),
	}

	lambda.Start(frontend.HandleRequest)
}
