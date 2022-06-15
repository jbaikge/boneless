package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/xid"
)

type Response struct {
	Id string `json:"id"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	r := Response{
		Id: xid.New().String(),
	}

	encoded, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return
	}

	response.Headers = map[string]string{
		"Content-Type": "application/json",
	}
	response.StatusCode = http.StatusOK
	response.Body = string(encoded)
	return
}

func main() {
	lambda.Start(HandleRequest)
}
