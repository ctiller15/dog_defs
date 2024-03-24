package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Welcome to dogdefs!",
	})
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html") // To load individual html files
	r.GET("/", homePage)

	r.Run(":8080")
}
