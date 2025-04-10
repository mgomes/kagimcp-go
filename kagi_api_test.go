package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestKagiSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/search" {
			t.Errorf("Expected path /api/v0/search, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bot testkey" {
			t.Errorf("Expected Authorization header 'Bot testkey', got '%s'", r.Header.Get("Authorization"))
		}
		query := r.URL.Query().Get("q")
		limit := r.URL.Query().Get("limit")

		if query != "test query" {
			t.Errorf("Expected query 'test query', got '%s'", query)
		}
		if limit != "5" {
			t.Errorf("Expected limit '5', got '%s'", limit)
		}

		resp := map[string]any{
			"data": []map[string]any{
				{
					"t":     0,
					"url":   "https://example.com/1",
					"title": "Test Result 1",
				},
				{
					"t":     0,
					"url":   "https://example.com/2",
					"title": "Test Result 2",
				},
				{
					"t":    1,
					"list": []string{"related1", "related2"},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	t.Log("Note: Direct testing of kagiSearch with httptest requires modifying the function to accept a base URL.")
}

func TestKagiSummarize(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v0/summarize" {
			t.Errorf("Expected path /api/v0/summarize, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bot testkey" {
			t.Errorf("Expected Authorization header 'Bot testkey', got '%s'", r.Header.Get("Authorization"))
		}
		targetURL := r.URL.Query().Get("url")
		engine := r.URL.Query().Get("engine")
		summaryType := r.URL.Query().Get("summary_type")

		if targetURL != "https://example.com/article" {
			t.Errorf("Expected url 'https://example.com/article', got '%s'", targetURL)
		}
		if engine != "cecil" {
			t.Errorf("Expected engine 'cecil', got '%s'", engine)
		}
		if summaryType != "summary" {
			t.Errorf("Expected summary_type 'summary', got '%s'", summaryType)
		}

		resp := map[string]any{
			"data": map[string]string{
				"output": "This is a test summary.",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	t.Log("Note: Direct testing of kagiSummarize with httptest requires modifying the function to accept a base URL.")
}
