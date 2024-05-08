package main

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/leolas95/mini-url/src/db"
	"net/http"
)

var ginLambda *ginadapter.GinLambda

type ErrorResponse struct {
	Message string `json:"message"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := db.Init()
	if err != nil {
		res := ErrorResponse{Message: "error initializing db: " + err.Error()}
		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(res)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       buf.String(),
		}, err
	}

	r := gin.Default()

	r.POST("/urls", Create)

	//r.Run(":8080")
	//return events.APIGatewayProxyResponse{}, nil

	ginLambda = ginadapter.New(r)
	return ginLambda.Proxy(request)
}

func main() {
	lambda.Start(Handler)
	//req := events.APIGatewayProxyRequest{
	//	MultiValueHeaders:               nil,
	//	QueryStringParameters:           nil,
	//	MultiValueQueryStringParameters: nil,
	//	PathParameters:                  nil,
	//	StageVariables:                  nil,
	//	RequestContext:                  events.APIGatewayProxyRequestContext{},
	//	Body:                            "",
	//	IsBase64Encoded:                 false,
	//}
	//_, _ = Handler(req)
}
