package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// SearchResult represents a single search result from Kagi
type SearchResult struct {
	Title       string
	URL         string
	Snippet     string
	PublishedAt string
}

// RelatedSearches represents related search terms from Kagi
type RelatedSearches struct {
	Terms []string
}

// SearchResponse contains all data returned from a search
type SearchResponse struct {
	Results         []SearchResult
	RelatedSearches []string
}

// kagiSearch performs the actual API call to Kagi search API
func (s *KagiServer) kagiSearch(ctx context.Context, query string, limit int) (SearchResponse, error) {
	apiKey := getAPIKey(ctx)
	if apiKey == "" {
		return SearchResponse{}, fmt.Errorf("kagi API key not found in context")
	}

	baseURL := "https://kagi.com/api/v0/search"
	params := make(url.Values)
	params.Add("q", query)
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}

	reqURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bot "+apiKey)

	client := &http.Client{Timeout: 29 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return SearchResponse{}, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	var apiResp struct {
		Data []struct {
			T         int      `json:"t"`
			URL       string   `json:"url,omitempty"`
			Title     string   `json:"title,omitempty"`
			Snippet   string   `json:"snippet,omitempty"`
			Published string   `json:"published,omitempty"`
			List      []string `json:"list,omitempty"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return SearchResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	var response SearchResponse

	for _, item := range apiResp.Data {
		if item.T == 0 {
			response.Results = append(response.Results, SearchResult{
				Title:       item.Title,
				URL:         item.URL,
				Snippet:     item.Snippet,
				PublishedAt: item.Published,
			})
		} else if item.T == 1 {
			response.RelatedSearches = item.List
		}
	}

	return response, nil
}

// kagiSummarize performs the actual API call to Kagi summarize API
func (s *KagiServer) kagiSummarize(ctx context.Context, targetURL, engine, summaryType string) (string, error) {
	apiKey := getAPIKey(ctx)
	if apiKey == "" {
		return "", fmt.Errorf("kagi API key not found in context")
	}

	baseURL := "https://kagi.com/api/v0/summarize"
	params := make(url.Values)
	params.Add("url", targetURL)
	params.Add("engine", engine)
	params.Add("summary_type", summaryType)

	reqURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bot "+apiKey)

	client := &http.Client{Timeout: 29 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, body)
	}

	var apiResp struct {
		Data struct {
			Output string `json:"output"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return apiResp.Data.Output, nil
}
