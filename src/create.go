package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

const baseURL = "http://miniurl.com/"

type Input struct {
	URL string `json:"url"`
}

type Result struct {
	ShortURL string `json:"short_url"`
}

func GenerateUniqueID(s string) (short string) {
	return uuid.New().String()
}

func Create(c *gin.Context) {
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// get unique url
	id := GenerateUniqueID(input.URL)
	short := c.Request.Host + c.Request.URL.Path + "/" + id
	fmt.Println(short)

	// store shorturl, longurl on db

	// get api base path

	// return shorturl to user

	res := Result{ShortURL: short}
	c.JSON(http.StatusOK, &res)
}
