package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"ctiller15/dog_defs/handlers"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"
	"github.com/supabase-community/supabase-go"
)

// CODE GO BRRRRRRRR

type wordResult struct {
	Id       int    `json:"word_id"`
	Word     string `json:"word"`
	Slug     string `json:"slug"`
	DefCount int    `json:"def_count"`
}

type paginationResponse struct {
	Value      string
	PageUrl    string
	IsCurrent  bool
	IsEllipsis bool
}

type DirectionalPageLink struct {
	Url     string
	Enabled bool
}

// Needs to be broken out and tested.
func createPagination(currentPage int, maxPage int, urlBase string) []paginationResponse {
	var pages []paginationResponse
	if currentPage > 2 && maxPage-currentPage > 2 {
		pages = []paginationResponse{
			{Value: "1", IsEllipsis: false, PageUrl: fmt.Sprintf("?page=%d%s", 1, urlBase)},
			{Value: "", IsEllipsis: true, PageUrl: ""},
			{Value: strconv.Itoa(currentPage - 1), IsEllipsis: false, PageUrl: fmt.Sprintf("?page=%d%s", currentPage-1, urlBase)},
			{Value: strconv.Itoa(currentPage), IsEllipsis: false, PageUrl: fmt.Sprintf("?page=%d%s", currentPage, urlBase), IsCurrent: true},
			{Value: strconv.Itoa(currentPage + 1), IsEllipsis: false, PageUrl: fmt.Sprintf("?page=%d%s", currentPage+1, urlBase)},
			{Value: "", IsEllipsis: true, PageUrl: ""},
			{Value: strconv.Itoa(maxPage), IsEllipsis: false, PageUrl: fmt.Sprintf("?page=%d%s", maxPage, urlBase)},
		}
	} else if currentPage > 2 {
		pages = []paginationResponse{
			{Value: "1", IsEllipsis: false, PageUrl: fmt.Sprintf("?page=%d%s", 1, urlBase)},
			{Value: "", IsEllipsis: true, PageUrl: ""},
		}

		for i := currentPage - 1; i < maxPage+1; i++ {
			pages = append(pages, paginationResponse{Value: strconv.Itoa(i), IsEllipsis: false, PageUrl: fmt.Sprintf("?page=%d%s", i, urlBase), IsCurrent: i == currentPage})
		}
	} else if maxPage-currentPage > 2 {
		for i := 1; i < currentPage+3; i++ {
			pages = append(pages, paginationResponse{Value: strconv.Itoa(i), IsEllipsis: false, PageUrl: fmt.Sprintf("?page=%d%s", i, urlBase), IsCurrent: i == currentPage})
		}

		pages = append(pages, paginationResponse{Value: "", IsEllipsis: true})
		pages = append(pages, paginationResponse{Value: strconv.Itoa(maxPage), IsEllipsis: false, PageUrl: fmt.Sprintf("?page=%d%s", maxPage, urlBase)})
	} else {
		// Loop up to current page. Loop up to max page.
		for i := 1; i < maxPage; i++ {
			pages = append(pages, paginationResponse{Value: strconv.Itoa(i), IsEllipsis: false, PageUrl: fmt.Sprintf("?page=%d%s", maxPage, urlBase), IsCurrent: i == currentPage})
		}
	}

	return pages
}

func wordsPage(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	pageSize := 40
	if err != nil {
		page = 1
	}

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

	data, count, err := client.
		From("words_list").
		Select("*", "exact", false).
		Range(pageSize*(page-1), pageSize*page, "").
		Execute()

	if err != nil {
		fmt.Println(err)
	}

	var wordsListResponse []wordResult

	json.Unmarshal(data, &wordsListResponse)

	fmt.Println(wordsListResponse)
	fmt.Println("^^^^^")

	maxPages := int(count/int64(pageSize)) + 1

	paginationResponse := createPagination(page, maxPages, "")

	previouspagenum := int(math.Max(float64(page-1), 1))
	nextpagenum := int(math.Min(float64(page+1), float64(maxPages)))

	c.HTML(http.StatusOK, "words.html", gin.H{
		"title":            "Words",
		"words":            &wordsListResponse,
		"count":            count,
		"pagination":       &paginationResponse,
		"previousPageLink": DirectionalPageLink{Url: "?page=" + strconv.Itoa(previouspagenum), Enabled: previouspagenum != page},
		"nextPageLink":     DirectionalPageLink{Url: "?page=" + strconv.Itoa(nextpagenum), Enabled: nextpagenum != page},
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
	Definition       template.HTML
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
			Definition:       template.HTML(model.Definition),
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
	word_slug := c.Param("word_slug")

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

	data, count, err := client.
		From("word_definitions").
		Select("*", "exact", false).
		Filter("slug", "eq", word_slug).
		Execute()

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
	p := bluemonday.UGCPolicy()
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
		"slug": slug.Make(strings.ToLower(newForm.Word)),
	}

	fmt.Println(saveWordForm)

	data, _, err := client.From("words").Insert(
		saveWordForm, true, "name", "representation", "",
	).Single().Execute()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("$$$")
	var wordResponse newDefinitionResponse
	json.Unmarshal(data, &wordResponse)

	word_id := wordResponse.Id

	// Now create the definition based on the word.
	saveDefinitionForm := map[string]interface{}{
		"text":           p.Sanitize(newForm.Definition),
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

	data, count, err := client.
		From("words_list").
		Select("*", "exact", false).
		Ilike("word", "%"+query+"%").
		Execute()

	if err != nil {
		fmt.Println(err)
	}

	var wordsListResponse []wordResult

	json.Unmarshal(data, &wordsListResponse)

	fmt.Println(&wordsListResponse)

	c.HTML(http.StatusOK, "search_page.html", gin.H{
		"title":   "Search results for " + query,
		"query":   query,
		"count":   count,
		"results": &wordsListResponse,
	})
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.StaticFile("/favicon.ico", "favicon.ico")
	r.LoadHTMLGlob("templates/*")

	handlers.SetupIndex(r)

	r.GET("/words", wordsPage)
	r.GET("/words/:word_slug", wordDefinitionPage)

	r.GET("/definitions/new", newDefinitionsPage)
	r.POST("/definitions/new", saveNewDefinition)

	r.GET("/search", searchPage)

	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file. May see unexpected behavior.")
	}

	r := setupRouter()

	r.Run(":8080")
}
