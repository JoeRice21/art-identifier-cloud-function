package main

import (
	"context"
	"fmt"
	"log"
	"os"

	p "github.com/JoeRice21/backend-serp-go"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	vertexAIProjectID := os.Getenv("VERTEX_PROJECT_ID")
	vertexAILocation := os.Getenv("VERTEX_PROJECT_LOCATION")

	ctx := p.NewTraceContext(context.Background())

	lensResult, err := p.MockSerpLensSearch()
	if err != nil {
		log.Fatal(err)
	}

	geminiResponse, err := p.GenerateArtworkResponse(ctx, lensResult, vertexAIProjectID, vertexAILocation)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", geminiResponse)
}