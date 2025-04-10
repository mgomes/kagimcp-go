package main

import (
	"context"
	"net/http"
)

// kagiAPIKeyKey is a context key for storing/retrieving the Kagi API key
type kagiAPIKeyKey struct{}

// withAPIKey adds the API key to the context for stdio transport
func (s *KagiServer) withAPIKey(ctx context.Context) context.Context {
	return context.WithValue(ctx, kagiAPIKeyKey{}, s.apiKey)
}

// withSSEAPIKey adds the API key to the context for SSE transport
func (s *KagiServer) withSSEAPIKey(ctx context.Context, r *http.Request) context.Context {
	apiKey := r.Header.Get("X-Kagi-API-Key")
	if apiKey == "" {
		apiKey = s.apiKey
	}
	return context.WithValue(ctx, kagiAPIKeyKey{}, apiKey)
}

// getAPIKey retrieves the Kagi API key from the context
func getAPIKey(ctx context.Context) string {
	if apiKey, ok := ctx.Value(kagiAPIKeyKey{}).(string); ok {
		return apiKey
	}
	return ""
}
