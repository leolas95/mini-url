package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	ShortURL string `json:"short_url"`
}

func Create(c *gin.Context) {
	// get unique url

	// store shorturl, longurl on db

	// get api base path

	// return shorturl to user

	res := Result{ShortURL: ""}
	c.JSON(http.StatusOK, &res)
}
