package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":    "DogDefs",
		"greeting": "Welcome to DogDefs!",
	})
}

func notFoundPage(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404_page.html", gin.H{
		"title": "Not Found",
	})
}

func SetupIndex(r *gin.Engine) *gin.Engine {
	r.GET("/", homePage)
	r.NoRoute(notFoundPage)
	return r
}
