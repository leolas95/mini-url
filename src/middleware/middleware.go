package middleware

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
)

// AddUrl appends our own Api gateway URL to the context, so we can retrieve it later
func AddUrl(request events.APIGatewayProxyRequest) gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme := c.Request.URL.Scheme
		host := request.Headers["Host"]
		stage := request.RequestContext.Stage
		requestPath := request.Path
		url := scheme + "://" + host + "/" + stage + requestPath
		c.Set("ApiGatewayURL", url)
		c.Next()
	}
}
