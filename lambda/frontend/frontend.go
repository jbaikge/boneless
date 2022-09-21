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
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jbaikge/boneless"
)

var (
	awsConfig aws.Config
	resources boneless.DynamoDBResources
)

type TemplateVars struct {
	Document boneless.Document
}

type Frontend struct {
	Repo boneless.Repository
}

func (frontend Frontend) HandleRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (response events.APIGatewayV2HTTPResponse, err error) {
	start := time.Now()
	documentService := boneless.NewDocumentService(frontend.Repo)
	document, byPathErr := documentService.ByPath(ctx, request.RawPath)
	if byPathErr != nil {
		response.StatusCode = http.StatusNotFound
		response.Body = "Document Not found!"
		return
	}

	vars := TemplateVars{
		Document: document,
	}

	buffer := new(bytes.Buffer)
	if compileErr := frontend.compileTemplates(ctx, vars, buffer); compileErr != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Body = fmt.Sprintf("template compilation error: %v", compileErr)
		return
	}

	response.StatusCode = http.StatusOK
	response.Headers = map[string]string{
		"Content-Type":   "text/html",
		"X-Handler-Time": time.Since(start).String(),
	}
	response.Body = buffer.String()
	return
}

func (frontend Frontend) compileTemplates(ctx context.Context, vars TemplateVars, w io.Writer) (err error) {
	templateService := boneless.NewTemplateService(frontend.Repo)
	filter := boneless.TemplateFilter{Range: boneless.Range{End: 1000}}
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
		if vars.Document.TemplateId == tmpl.Id {
			name = tmpl.Id
		}
		if _, err = t.New(name).Parse(tmpl.Body); err != nil {
			return fmt.Errorf("parsing %s: %w", tmpl.Name, err)
		}
	}

	return t.ExecuteTemplate(w, vars.Document.TemplateId, vars)
}

func (frontend Frontend) decodeFilter(s string) (filter boneless.DocumentFilter, err error) {
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
	classService := boneless.NewClassService(frontend.Repo)
	classes, err := classService.All(context.Background())
	if err != nil {
		return
	}

	classNameMap := make(map[string]string)
	for _, class := range classes {
		classNameMap[class.Name] = class.Id
	}

	return template.FuncMap{
		"get_document": func(id string) (doc boneless.Document, err error) {
			return boneless.NewDocumentService(frontend.Repo).ById(context.Background(), id)
		},
		"list_documents": func(className string, args string) (docs []boneless.Document, err error) {
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

			documentService := boneless.NewDocumentService(frontend.Repo)
			docs, _, err = documentService.List(context.Background(), filter)
			return
		},
		"many_documents": func(ids []string) (docs []boneless.Document, err error) {
			docs = make([]boneless.Document, 0, len(ids))
			documentService := boneless.NewDocumentService(frontend.Repo)
			for _, id := range ids {
				doc, err := documentService.ById(context.Background(), id)
				if err != nil {
					return nil, err
				}
				docs = append(docs, doc)
			}
			return
		},
		"child_documents": func(className string, parentId string) (docs []boneless.Document, err error) {
			id, found := classNameMap[className]
			if !found {
				err = fmt.Errorf("invalid class name: %s", className)
				return
			}

			filter := boneless.DocumentFilter{
				ClassId:  id,
				ParentId: parentId,
				Range:    boneless.Range{End: 100},
			}
			docs, _, err = boneless.NewDocumentService(frontend.Repo).List(context.Background(), filter)
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
		Repo: boneless.NewDynamoDBRepository(awsConfig, resources),
	}

	lambda.Start(frontend.HandleRequest)
}
