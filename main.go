package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":    "DogDefs",
		"greeting": "Welcome to dogdefs!",
	})
}

func saveNewDefinition(c *gin.Context) {
	fmt.Println("Saving definition!")
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html") // To load individual html files
	r.GET("/", homePage)

	r.POST("/definitions/new", saveNewDefinition)

	r.Run(":8080")
}
