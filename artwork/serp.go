package artwork

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func MockSerpLensSearch() (*LensResponse, error) {
	data, err := os.ReadFile("mock_about_this_image_response.json")
	if err != nil {
		return nil, err
	}

	var result LensResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func SerpLensSearch(ctx context.Context, imageURL string, serpAPIKey string) (*LensResponse, error) {
	logger := GetLogger(ctx, "helper::serpLensSearch")
	logger.Info("starting serp lens search", "imageURL", imageURL)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	params := url.Values{}
	params.Set("engine", "google_lens")
	params.Set("api_key", serpAPIKey)
	params.Set("url", imageURL)
	params.Set("hl", "en")
	params.Set("country", "gb")
	params.Set("type", "about_this_image")

	endpoint := "https://serpapi.com/search?" + params.Encode()

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		logger.Error("failed to build request", "error", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("serp api call failed", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("unexpected response status", "status", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	logger.Info("serp api call successful", "status", resp.StatusCode)

	var result LensResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("failed to decode response", "error", err)
		return nil, err
	}
	logger.Info("serp response decoded", "sections", len(result.AboutThisImage.Sections))

	return &result, nil
}