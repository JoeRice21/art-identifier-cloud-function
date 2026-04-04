package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	vertexAIProjectID := os.Getenv("VERTEX_PROJECT_ID")
	vertexAILocation := os.Getenv("VERTEX_PROJECT_LOCATION")

	lensSearchResult, err := mockSerpLensSearch()
	if err != nil {
		log.Fatal(err)
	}

	geminiResponse, err := generateArtworkResponse(lensSearchResult, vertexAIProjectID, vertexAILocation)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", geminiResponse)
}