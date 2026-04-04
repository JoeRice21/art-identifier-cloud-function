package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func mockSerpLensSearch() (*LensResponse, error) {
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

func serpLensSearch(imageURL string, serpAPIKey string) (*LensResponse, error) {
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
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result LensResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}