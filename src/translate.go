package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/leolas95/mini-url/src/db"
	"net/http"
)

func Translate(c *gin.Context) {
	shortUrl := c.Param("url")

	// check url on db
	out, err := db.DB.GetItem(c, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: shortUrl},
		},
		TableName:            aws.String("urls"),
		ProjectionExpression: aws.String("id,long_url"),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving item: " + err.Error()})
		return
	}
	item := out.Item

	if value, ok := item["long_url"]; !ok || value == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found - check that you entered the correct short url"})
		return
	}

	longUrl := item["long_url"].(*types.AttributeValueMemberS).Value

	c.Header("location", longUrl)
	c.JSON(http.StatusTemporaryRedirect, "Redirecting you...")
}
