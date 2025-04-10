package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// handleKagiSearch implements the Kagi search API tool
func (s *KagiServer) handleKagiSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid query parameter")
	}

	limit := 5
	if limitParam, ok := request.Params.Arguments["limit"].(float64); ok {
		limit = int(limitParam)
	}

	searchResponse, err := s.kagiSearch(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("kagi search failed: %w", err)
	}

	resultText := fmt.Sprintf("Search results for '%s':\n\n", query)
	for i, result := range searchResponse.Results {
		resultText += fmt.Sprintf("%d. %s\n   URL: %s\n", i+1, result.Title, result.URL)
		if result.PublishedAt != "" {
			resultText += fmt.Sprintf("   Published: %s\n", result.PublishedAt)
		}
		resultText += fmt.Sprintf("   %s\n\n", result.Snippet)
	}

	if len(searchResponse.RelatedSearches) > 0 {
		resultText += "Related searches:\n"
		for i, term := range searchResponse.RelatedSearches {
			resultText += fmt.Sprintf("%d. %s\n", i+1, term)
		}
	}

	return mcp.NewToolResultText(resultText), nil
}

// handleKagiSummarize implements the Kagi summarize API tool
func (s *KagiServer) handleKagiSummarize(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, ok := request.Params.Arguments["url"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid url parameter")
	}

	engine := "agnes"
	if engineParam, ok := request.Params.Arguments["engine"].(string); ok {
		engine = engineParam
	}

	summaryType := "summary"
	if summaryTypeParam, ok := request.Params.Arguments["summary_type"].(string); ok {
		summaryType = summaryTypeParam
	}

	summary, err := s.kagiSummarize(ctx, url, engine, summaryType)
	if err != nil {
		return nil, fmt.Errorf("kagi summarize failed: %w", err)
	}

	resultText := fmt.Sprintf("Summary of %s:\n\n%s", url, summary)
	return mcp.NewToolResultText(resultText), nil
}
