package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

// CODE GO BRRRRRRRR

func homePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":    "DogDefs",
		"greeting": "Welcome to Dogdefs!",
	})
}

type wordResult struct {
	Id       int    `json:"word_id"`
	Word     string `json:"word"`
	DefCount int    `json:"def_count"`
}

func wordsPage(c *gin.Context) {
	fmt.Println("Fetching all approved words...")

	// TODO: Break data retrieval into separate module.
	apiUrl := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	fmt.Println(apiUrl)
	fmt.Println(apiKey)

	client, err := supabase.NewClient(
		apiUrl,
		apiKey,
		nil)

	if err != nil {
		fmt.Println("cannot initialize client", err)
	}

	data, count, err := client.From("words_list").Select("*", "exact", false).Execute()

	if err != nil {
		fmt.Println(err)
	}

	var wordsListResponse []wordResult

	json.Unmarshal(data, &wordsListResponse)

	c.HTML(http.StatusOK, "words.html", gin.H{
		"title": "Words",
		"words": &wordsListResponse,
		"count": count,
	})
}

type wordDefinitionListResult struct {
	WordId        int    `json:"id"`
	Word          string `json:"name"`
	Definition    string `json:"text"`
	PartOfSpeech  string `json:"part_of_speech"`
	Reference     string `json:"reference"`
	ReferenceLink string `json:"reference_link"`
}

type definitionViewModel struct {
	WordId           int
	Word             string
	Definition       string
	PartOfSpeech     string
	Reference        string
	HasReference     bool
	ReferenceLink    string
	HasReferenceLink bool
}

func mapToDefinitionViewModel(inModel []wordDefinitionListResult) []definitionViewModel {
	output := make([]definitionViewModel, len(inModel))

	for i, model := range inModel {
		viewModel := definitionViewModel{
			WordId:           model.WordId,
			Word:             model.Word,
			Definition:       model.Definition,
			PartOfSpeech:     model.PartOfSpeech,
			Reference:        model.Reference,
			HasReference:     len(model.Reference) > 0,
			ReferenceLink:    model.ReferenceLink,
			HasReferenceLink: len(model.ReferenceLink) > 0,
		}

		output[i] = viewModel
	}

	return output
}

func wordDefinitionPage(c *gin.Context) {
	fmt.Println("Getting definitions for word...")
	word_id := c.Param("word_id")

	// TODO: Break data retrieval into separate module.
	apiUrl := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")

	client, err := supabase.NewClient(
		apiUrl,
		apiKey,
		nil)

	if err != nil {
		fmt.Println("cannot initialize client", err)
	}

	var listResult []wordDefinitionListResult

	data, count, err := client.From("word_definitions").Select("*", "exact", false).Filter("id", "eq", word_id).Execute()

	// If count is somehow less than 1 we want to return a 404.

	json.Unmarshal(data, &listResult)

	id := listResult[0].WordId
	word_name := listResult[0].Word

	listView := mapToDefinitionViewModel(listResult)

	c.HTML(http.StatusOK, "word_definitions.html", gin.H{
		"word_id":          id,
		"word":             word_name,
		"definition_count": count,
		"definitions":      &listView,
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
	Link         string `form:"link"`
}

type newDefinitionResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func saveNewDefinition(c *gin.Context) {
	// TODO: Break data retrieval into separate module.
	apiUrl := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")

	client, err := supabase.NewClient(
		apiUrl,
		apiKey,
		nil)

	if err != nil {
		fmt.Println("cannot initialize client", err)
	}

	var newForm newDefinitionForm
	c.ShouldBind(&newForm)

	saveWordForm := map[string]interface{}{
		"name": strings.ToLower(newForm.Word),
	}

	data, _, err := client.From("words").Insert(
		saveWordForm, true, "name", "representation", "",
	).Single().Execute()

	if err != nil {
		fmt.Println(err)
	}

	var wordResponse newDefinitionResponse
	json.Unmarshal(data, &wordResponse)

	fmt.Println(wordResponse)

	word_id := wordResponse.Id
	fmt.Println(word_id)

	// Now create the definition based on the word.
	saveDefinitionForm := map[string]interface{}{
		"text":           newForm.Definition,
		"approved":       false,
		"word_id":        word_id,
		"part_of_speech": newForm.PartOfSpeech,
		"reference":      newForm.Reference,
		"reference_link": newForm.Link,
	}

	_, _, err = client.From("definitions").Insert(
		saveDefinitionForm, false, "", "", "",
	).Single().Execute()

	if err != nil {
		fmt.Println(err)
	}

	c.HTML(http.StatusOK, "definition_saved.html", gin.H{
		"title": "Definition Saved",
		"word":  newForm.Word,
	})
}

func searchPage(c *gin.Context) {
	query := c.Query("query")
	fmt.Println("Searching words!")
	fmt.Println(query)

	// TODO: Break data retrieval into separate module.
	apiUrl := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")

	client, err := supabase.NewClient(
		apiUrl,
		apiKey,
		nil)

	if err != nil {
		fmt.Println("cannot initialize client", err)
	}

	data, count, err := client.From("words_list").Select("*", "exact", false).Ilike("word", "%"+query+"%").Execute()

	if err != nil {
		fmt.Println(err)
	}

	var wordsListResponse []wordResult

	json.Unmarshal(data, &wordsListResponse)

	c.HTML(http.StatusOK, "search_page.html", gin.H{
		"title":   "Search results for " + query,
		"query":   query,
		"count":   count,
		"results": &wordsListResponse,
	})
}

func notFoundPage(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404_page.html", gin.H{
		"title": "Not Found",
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file. May see unexpected behavior.")
	}
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", homePage)
	r.GET("/words", wordsPage)
	r.GET("/words/:word_id", wordDefinitionPage)

	r.GET("/definitions/new", newDefinitionsPage)
	r.POST("/definitions/new", saveNewDefinition)

	r.GET("/search", searchPage)

	r.NoRoute(notFoundPage)

	r.Run(":8080")
}
