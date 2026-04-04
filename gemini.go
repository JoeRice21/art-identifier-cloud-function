package p

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/genai"
)

func GenerateArtworkResponse(ctx context.Context, lensResponse *LensResponse, vertexAIProjectID string, vertexAILocation string) (*ArtworkResponse, error) {
	logger := GetLogger(ctx, "helper::generateArtworkResponse")
	logger.Info("starting artwork identification")

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  vertexAIProjectID,
		Location: vertexAILocation,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		logger.Error("failed to create genai client", "error", err)
		return nil, err
	}
	logger.Info("genai client created")

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

	logger.Info("sending prompt to gemini", "model", "gemini-2.5-flash")
	resp, err := client.Models.GenerateContent(ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		logger.Error("gemini generate content failed", "error", err)
		return nil, err
	}
	logger.Info("gemini response received")

	var artworkResponse ArtworkResponse
	if err := json.Unmarshal([]byte(resp.Text()), &artworkResponse); err != nil {
		logger.Error("failed to unmarshal gemini response", "error", err, "rawResponse", resp.Text())
		return nil, err
	}

	logger.Info("artwork identified", "artist", artworkResponse.Artist, "title", artworkResponse.ArtworkTitle)
	return &artworkResponse, nil
}