package p

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("ArtworkHandler", ArtworkHandler)
}

func ArtworkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := NewTraceContext(r.Context())
	logger := GetLogger(ctx, "controller::ArtworkHandler")

	imageURL := r.URL.Query().Get("imageURL")
	logger.Info("received request", "imageURL", imageURL)

	if imageURL == "" {
		logger.Error("missing imageURL parameter")
		http.Error(w, "missing imageURL query parameter", http.StatusBadRequest)
		return
	}

	serpAPIKey := os.Getenv("SERP_KEY")
	vertexAIProjectID := os.Getenv("VERTEX_PROJECT_ID")
	vertexAILocation := os.Getenv("VERTEX_PROJECT_LOCATION")

	lensResult, err := SerpLensSearch(ctx, imageURL, serpAPIKey)
	if err != nil {
		logger.Error("serp lens search failed", "error", err)
		http.Error(w, "lens search failed", http.StatusInternalServerError)
		return
	}
	logger.Info("serp lens search successful", "sections", len(lensResult.AboutThisImage.Sections))

	artworkResponse, err := GenerateArtworkResponse(ctx, lensResult, vertexAIProjectID, vertexAILocation)
	if err != nil {
		logger.Error("gemini failed", "error", err)
		http.Error(w, "gemini failed", http.StatusInternalServerError)
		return
	}

	logger.Info("request complete", "artist", artworkResponse.Artist, "title", artworkResponse.ArtworkTitle)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artworkResponse)
}