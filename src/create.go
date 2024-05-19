package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/leolas95/mini-url/src/db"
	"github.com/leolas95/mini-url/src/util"
	"net/http"
	"net/url"
)

const (
	genUrlIdServiceUrl = "https://kajlwtj945.execute-api.us-east-1.amazonaws.com/prod/"
)

type CreateInput struct {
	URL string `json:"url"`
}

type CreateResponse struct {
	ShortURL string `json:"short_url"`
}

// Types for calling the gen-url-id service

type GenUrlRequestBody struct {
	Url string `json:"url"`
}

type GenUrlResponseBody struct {
	Hash string `json:"hash"`
}

func saveItemOnDB(id, longUrl string) error {
	item := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: id},
			"long_url": &types.AttributeValueMemberS{Value: longUrl},
		},
		TableName: aws.String("urls"),
	}
	_, err := db.DB.PutItem(context.TODO(), &item)
	return err
}

func Create(c *gin.Context) {
	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endpointUrl, err := url.JoinPath(genUrlIdServiceUrl, "/urls")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error building url endpoint for hash generation service: " + err.Error()})
		return
	}

	// get unique id
	body := GenUrlRequestBody{Url: input.URL}
	response, err := util.MakePostRequest(endpointUrl, "application/json", body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error making POST request to hash generation service:" + err.Error()})
		return
	}
	genUrlResponse, ok := response.(*GenUrlResponseBody)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing hash generation service response"})
		return
	}
	id := genUrlResponse.Hash

	// get our own api base path
	short := c.GetString("ApiGatewayURL") + "/" + id

	// save data on db
	err = saveItemOnDB(id, input.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error saving item on DB: " + err.Error()})
		return
	}

	// return shorturl to user
	res := CreateResponse{ShortURL: short}
	c.JSON(http.StatusOK, &res)
}
