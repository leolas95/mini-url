package main

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/leolas95/mini-url/src/db"
	"io"
	"net/http"
	"net/url"
)

const (
	genUrlIdServiceUrl = "https://kajlwtj945.execute-api.us-east-1.amazonaws.com/prod/"
)

type CreateInput struct {
	URL string `json:"url"`
}

type CreateResult struct {
	ShortURL string `json:"short_url"`
}

// Types for calling the gen-url-id service

type GenUrlRequestBody struct {
	Url string `json:"url"`
}

type GenUrlResponseBody struct {
	Hash string `json:"hash"`
}

func makePostRequest(url, contentType string, body GenUrlRequestBody) (*GenUrlResponseBody, error) {
	marshaled, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	postBody := bytes.NewBuffer(marshaled)
	resp, err := http.Post(url, contentType, postBody)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var serviceResponse GenUrlResponseBody
	err = json.Unmarshal(bytesResp, &serviceResponse)
	if err != nil {
		return nil, err
	}

	return &serviceResponse, nil
}

func Create(c *gin.Context) {
	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get unique id
	endpointUrl, err := url.JoinPath(genUrlIdServiceUrl, "/urls")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error building url endpoint for hash generation service: " + err.Error()})
		return
	}

	body := GenUrlRequestBody{Url: input.URL}
	genUrlResponse, err := makePostRequest(endpointUrl, "application/json", body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error making POST request to hash generation service:" + err.Error()})
		return
	}
	id := genUrlResponse.Hash

	// get this service api base path
	short := c.GetString("ApiGatewayURL") + "/" + id

	// save data on db
	item := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: id},
			"long_url": &types.AttributeValueMemberS{Value: input.URL},
		},
		TableName: aws.String("urls"),
	}
	_, err = db.DB.PutItem(c, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// return shorturl to user
	res := CreateResult{ShortURL: short}
	c.JSON(http.StatusOK, &res)
}
