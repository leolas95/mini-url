package main

import (
	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r := gin.Default()

	r.POST("/urls", Create)

	r.Run(":8080")
	return events.APIGatewayProxyResponse{}, nil

	//ginLambda = ginadapter.New(r)
	//return ginLambda.Proxy(request)
}

func main() {
	//lambda.Start(Handler)
	req := events.APIGatewayProxyRequest{
		MultiValueHeaders:               nil,
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		PathParameters:                  nil,
		StageVariables:                  nil,
		RequestContext:                  events.APIGatewayProxyRequestContext{},
		Body:                            "",
		IsBase64Encoded:                 false,
	}
	_, _ = Handler(req)
}
