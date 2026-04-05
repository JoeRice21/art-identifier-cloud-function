package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func OrchestratorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := NewTraceContext(r.Context())
	logger := GetLogger(ctx, "controller::OrchestratorHandler")

	// step 1 — forward the image to the upload function
	logger.Info("calling upload function")
	uploadResp, err := callUploadFunction(ctx, r)
	if err != nil {
		logger.Error("upload function failed", "error", err)
		http.Error(w, "failed to upload image", http.StatusInternalServerError)
		return
	}
	logger.Info("upload successful", "signedURL", uploadResp.SignedURL)

	// step 2 — pass the signed URL to the artwork function
	logger.Info("calling artwork function")
	artworkResp, err := callArtworkFunction(ctx, uploadResp.SignedURL)
	if err != nil {
		logger.Error("artwork function failed", "error", err)
		http.Error(w, "failed to identify artwork", http.StatusInternalServerError)
		return
	}
	logger.Info("artwork identified", "artist", artworkResp.Artist, "title", artworkResp.ArtworkTitle)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artworkResp)
}

type uploadResponse struct {
	SignedURL string `json:"signedURL"`
}

func callUploadFunction(ctx context.Context, r *http.Request) (*uploadResponse, error) {
	uploadURL := os.Getenv("UPLOAD_FUNCTION_URL")

	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to build upload request: %w", err)
	}
	req.Header = r.Header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload function returned %s", resp.Status)
	}

	var uploadResp uploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, err
	}

	return &uploadResp, nil
}

func callArtworkFunction(ctx context.Context, imageURL string) (*ArtworkResponse, error) {
	artworkURL := os.Getenv("ARTWORK_FUNCTION_URL")

	req, err := http.NewRequestWithContext(ctx, "GET", artworkURL+"?imageURL="+url.QueryEscape(imageURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build artwork request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("artwork function returned %s", resp.Status)
	}

	var artworkResp ArtworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&artworkResp); err != nil {
		return nil, err
	}

	return &artworkResp, nil
}

func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}