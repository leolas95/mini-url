package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ginLambda *ginadapter.GinLambda

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "Hello World") })

	ginLambda = ginadapter.New(r)
	return ginLambda.Proxy(request)
}

func main() {
	lambda.Start(Handler)
}
