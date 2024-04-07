package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":    "DogDefs",
		"greeting": "Welcome to Dogdefs!",
	})
}

func NotFoundPage(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404_page.html", gin.H{
		"title": "Not Found",
	})
}
