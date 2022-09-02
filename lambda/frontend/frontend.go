package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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

type TemplateVars struct {
	Document gocms.Document
}

type Frontend struct {
	Repo gocms.Repository
	vars TemplateVars
}

func (frontend Frontend) HandleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	documentService := gocms.NewDocumentService(frontend.Repo)
	document, byPathErr := documentService.ByPath(ctx, request.RawPath)
	if byPathErr != nil {
		response.StatusCode = http.StatusNotFound
		response.Body = "Document Not found!"
		return
	}

	frontend.vars.Document = document

	buffer := new(bytes.Buffer)
	if compileErr := frontend.compileTemplates(ctx, buffer); compileErr != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = fmt.Sprintf("template compilation error: %v", compileErr)
		return
	}

	response.StatusCode = http.StatusOK
	response.Headers = map[string]string{
		"Content-Type": "text/html",
	}
	response.Body = buffer.String()
	return
}

func (frontend Frontend) compileTemplates(ctx context.Context, w io.Writer) (err error) {
	templateService := gocms.NewTemplateService(frontend.Repo)
	filter := gocms.TemplateFilter{Range: gocms.Range{End: 1000}}
	templates, _, err := templateService.List(ctx, filter)
	if err != nil {
		return fmt.Errorf("fetch templates: %w", err)
	}

	funcs, err := frontend.funcMap()
	if err != nil {
		return fmt.Errorf("building func map: %w", err)
	}

	t := template.New("").Funcs(funcs)
	for _, tmpl := range templates {
		name := tmpl.Name
		if frontend.vars.Document.TemplateId == tmpl.Id {
			name = tmpl.Id
		}
		if _, err = t.New(name).Parse(tmpl.Body); err != nil {
			return fmt.Errorf("parsing %s: %w", tmpl.Name, err)
		}
	}

	return t.ExecuteTemplate(w, frontend.vars.Document.TemplateId, frontend.vars)
}

func (frontend Frontend) decodeFilter(s string) (filter gocms.DocumentFilter, err error) {
	for _, arg := range strings.Split(s, ";") {
		key, value, found := strings.Cut(arg, ":")
		if !found {
			return filter, fmt.Errorf("no value found for key: %s", key)
		}
		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.TrimSpace(value)
		switch key {
		case "range":
			lower, upper, found := strings.Cut(value, "-")
			if !found {
				upper = lower
				lower = "0"
			}
			if filter.Range.Start, err = strconv.Atoi(lower); err != nil {
				return filter, fmt.Errorf("converting range start: %w", err)
			}
			if filter.Range.End, err = strconv.Atoi(upper); err != nil {
				return filter, fmt.Errorf("converting range end: %w", err)
			}
		case "sort":
			field, dir, found := strings.Cut(value, ",")
			if !found {
				dir = "ASC"
			}
			filter.Sort.Field = strings.TrimSpace(field)
			filter.Sort.Direction = strings.TrimSpace(dir)
		case "parent":
			filter.ParentId = value
		}
	}

	return
}

func (frontend Frontend) funcMap() (funcs template.FuncMap, err error) {
	classService := gocms.NewClassService(frontend.Repo)
	classes, err := classService.All(context.Background())
	if err != nil {
		return
	}

	classNameMap := make(map[string]string)
	for _, class := range classes {
		classNameMap[class.Name] = class.Id
	}

	return template.FuncMap{
		"get_document": func(id string) (doc gocms.Document, err error) {
			return gocms.NewDocumentService(frontend.Repo).ById(context.Background(), id)
		},
		"list_documents": func(className string, args string) (docs []gocms.Document, err error) {
			id, found := classNameMap[className]
			if !found {
				err = fmt.Errorf("invalid class name: %s", className)
				return
			}

			filter, err := frontend.decodeFilter(args)
			if err != nil {
				return
			}
			filter.ClassId = id

			documentService := gocms.NewDocumentService(frontend.Repo)
			docs, _, err = documentService.List(context.Background(), filter)
			return
		},
		"many_documents": func(ids []string) (docs []gocms.Document, err error) {
			docs = make([]gocms.Document, 0, len(ids))
			documentService := gocms.NewDocumentService(frontend.Repo)
			for _, id := range ids {
				doc, err := documentService.ById(context.Background(), id)
				if err != nil {
					return nil, err
				}
				docs = append(docs, doc)
			}
			return
		},
		"split": strings.Fields,
	}, nil
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
