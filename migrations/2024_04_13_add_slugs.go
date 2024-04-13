package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type wordModel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file. May see unexpected behavior.")
	}

	fmt.Println("Running migration!")

	apiUrl := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	fmt.Println(apiUrl)

	client, err := supabase.NewClient(
		apiUrl,
		apiKey,
		nil)

	if err != nil {
		fmt.Println("cannot initialize client", err)
	}

	data, count, err := client.
		From("words").
		Select("*", "exact", false).
		Is("slug", "null").
		Execute()

	var unSluggedWords []wordModel
	json.Unmarshal(data, &unSluggedWords)

	for _, word := range unSluggedWords {
		fmt.Println(word)
		new_slug := slug.Make(word.Name)
		fmt.Println(new_slug)
		updateSlugBody := map[string]interface{}{
			"slug": new_slug,
		}
		data, _, err := client.
			From("words").
			Update(updateSlugBody, "representation", "").
			Filter("id", "eq", strconv.Itoa(word.Id)).
			Execute()

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(data)
	}

	// fmt.Println(data)
	fmt.Println(count)
}
