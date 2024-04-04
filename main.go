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

func wordsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "words.html", gin.H{
		"title": "Words",
	})
}

func newDefinitionsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "new_definition.html", gin.H{
		"title": "New Definition",
	})
}

type newDefinitionForm struct {
	Word         string `form:"word"`
	PartOfSpeech string `form:"pos"`
	Definition   string `form:"definition"`
	Reference    string `form:"reference"`
}

func saveNewDefinition(c *gin.Context) {
	fmt.Println("Saving definition!")
	var newForm newDefinitionForm
	c.ShouldBind(&newForm)
	c.HTML(http.StatusOK, "definition_saved.html", gin.H{
		"title": "Definition Saved",
		"word":  newForm.Word,
	})
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html") // To load individual html files
	r.GET("/", homePage)
	r.GET("/words", wordsPage)

	r.GET("/definitions/new", newDefinitionsPage)
	r.POST("/definitions/new", saveNewDefinition)

	r.Run(":8080")
}
