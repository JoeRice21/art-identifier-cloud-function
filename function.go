package p

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("ArtworkHandler", ArtworkHandler)
}

func ArtworkHandler(w http.ResponseWriter, r *http.Request) {
	imageURL := r.URL.Query().Get("imageURL")
	if imageURL == "" {
		http.Error(w, "missing imageURL query parameter", http.StatusBadRequest)
		return
	}

	serpAPIKey := os.Getenv("SERP_KEY")
	vertexAIProjectID := os.Getenv("VERTEX_PROJECT_ID")
	vertexAILocation := os.Getenv("VERTEX_PROJECT_LOCATION")

	lensResult, err := SerpLensSearch(imageURL, serpAPIKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("lens search failed: %s", err), http.StatusInternalServerError)
		return
	}

	artworkResponse, err := GenerateArtworkResponse(lensResult, vertexAIProjectID, vertexAILocation)
	if err != nil {
		http.Error(w, fmt.Sprintf("gemini failed: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artworkResponse)
}