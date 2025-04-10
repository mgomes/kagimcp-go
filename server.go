package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// KagiServer represents our MCP server with Kagi API support
type KagiServer struct {
	apiKey    string
	mcpServer *server.MCPServer
}

// NewKagiServer creates a new Kagi MCP server
func NewKagiServer(apiKey string) *KagiServer {
	mcpServer := server.NewMCPServer(
		"KagiMCP",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithInstructions("This server provides access to Kagi Search and Summarizer APIs."),
	)

	kagiServer := &KagiServer{
		apiKey:    apiKey,
		mcpServer: mcpServer,
	}

	kagiServer.registerTools()

	return kagiServer
}

func (s *KagiServer) registerTools() {
	s.mcpServer.AddTool(
		mcp.NewTool("kagi_search",
			mcp.WithDescription("Search the web using Kagi"),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("The search query string"),
			),
			mcp.WithNumber("limit",
				mcp.Description("Maximum number of results (1-10)"),
				mcp.Min(1),
				mcp.Max(10),
				mcp.DefaultNumber(5),
			),
		),
		s.handleKagiSearch,
	)

	s.mcpServer.AddTool(
		mcp.NewTool("kagi_summarize",
			mcp.WithDescription("Summarize a webpage using Kagi's Universal Summarizer API"),
			mcp.WithString("url",
				mcp.Required(),
				mcp.Description("URL of the webpage to summarize"),
			),
			mcp.WithString("engine",
				mcp.Description("Summarization engine to use (cecil, agnes, muriel)"),
				mcp.Enum("cecil", "agnes", "muriel"),
				mcp.DefaultString("agnes"),
			),
			mcp.WithString("summary_type",
				mcp.Description("Summary types are control the structure of the summary output (summary, takeaway)"),
				mcp.Enum("summary", "takeaway"),
				mcp.DefaultString("summary"),
			),
		),
		s.handleKagiSummarize,
	)
}

// ServeStdio starts the server using stdio transport
func (s *KagiServer) ServeStdio() error {
	return server.ServeStdio(s.mcpServer, server.WithStdioContextFunc(s.withAPIKey))
}

// ServeSSE starts the server using SSE transport
func (s *KagiServer) ServeSSE(addr string) error {
	sseServer := server.NewSSEServer(s.mcpServer,
		server.WithBaseURL(fmt.Sprintf("http://%s", addr)),
		server.WithSSEContextFunc(s.withSSEAPIKey),
	)
	return sseServer.Start(addr)
}
