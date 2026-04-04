package main

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/genai"
)

func generateArtworkResponse(lensResponse *LensResponse, vertexAIProjectID string, vertexAILocation string) (*ArtworkResponse, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  vertexAIProjectID,
		Location: vertexAILocation,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		return nil, err
	}

	geminiResponseSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"artist": map[string]any{
				"type":        "string",
				"description": "Name of the artist",
			},
			"artworkTitle": map[string]any{
				"type":        "string",
				"description": "Title of the artwork",
			},
			"dateOfCreation": map[string]any{
				"type":        "string",
				"description": "Date the artwork was created",
			},
			"artworkDescription": map[string]any{
				"type":        "string",
				"description": "2-3 sentence brief description of the artwork",
			},
			"wikipediaLink": map[string]any{
				"type":        "string",
				"description": "Link to the artworks wikipedia artwork",
			},
		},
		"required": []string{"artist", "artworkTitle", "dateOfCreation", "artworkDescription", "wikipediaLink"},
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: geminiResponseSchema,
	}

	prompt := fmt.Sprintf(`Based on these search results from a reverse image search, identify the artwork, artist and wikipedia entry for said artwork:

    %s
	`, lensResponse)

	resp, err := client.Models.GenerateContent(ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return nil, err
	}

	var artworkResponse ArtworkResponse
	if err := json.Unmarshal([]byte(resp.Text()), &artworkResponse); err != nil {
		return nil, err
	}

	return &artworkResponse, nil
}