package main

import (
	"testing"
)

func TestNewKagiServer(t *testing.T) {
	apiKey := "test-api-key"
	kagiServer := NewKagiServer(apiKey)

	if kagiServer == nil {
		t.Fatal("NewKagiServer should return a non-nil KagiServer")
	}

	if kagiServer.apiKey != apiKey {
		t.Errorf("Expected apiKey to be %q, got %q", apiKey, kagiServer.apiKey)
	}

	if kagiServer.mcpServer == nil {
		t.Fatal("mcpServer should be initialized")
	}

	// We can't directly test the internal configuration of mcpServer
	// as these are unexported fields.
}
