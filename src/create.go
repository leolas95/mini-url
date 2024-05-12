package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leolas95/mini-url/src/db"
	"net/http"
)

type CreateInput struct {
	URL string `json:"url"`
}

type CreateResult struct {
	ShortURL string `json:"short_url"`
}

func GenerateUniqueID(s string) (short string) {
	return uuid.New().String()
}

func Create(c *gin.Context) {
	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// get unique id
	id := GenerateUniqueID(input.URL)

	// get api base path
	short := c.GetString("ApiGatewayURL") + "/" + id

	// store shorturl, longurl on db
	item := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: id},
			"long_url": &types.AttributeValueMemberS{Value: input.URL},
		},
		TableName: aws.String("urls"),
	}
	_, err := db.DB.PutItem(c, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// return shorturl to user
	res := CreateResult{ShortURL: short}
	c.JSON(http.StatusOK, &res)
}
