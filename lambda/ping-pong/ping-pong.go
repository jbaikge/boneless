package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/xid"
)

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
	fmt.Printf("Processing request %s\n", request.RequestContext.RequestID)

	response.Headers = map[string]string{
		"X-XID": xid.New().String(),
	}
	response.StatusCode = http.StatusOK
	response.Body = request.Body
	return
}

func main() {
	lambda.Start(HandleRequest)
}
