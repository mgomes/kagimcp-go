# Kagi MCP Server

A Model Context Protocol (MCP) server that integrates with Kagi's search and summarizer APIs. This server allows Large Language Models (LLMs) to search the web and summarize web pages using Kagi's high-quality search and AI capabilities.

## Features

- 🔍 **Kagi Search**: Search the web with Kagi's privacy-focused search engine
- 📝 **Kagi Summarizer**: Summarize web pages using Kagi's FastGPT summarization APIs
- 🔄 **Multiple Transports**: Support for both stdio and Server-Sent Events (SSE) protocols
- 🔑 **API Key Management**: Flexible options for providing Kagi API keys

## Prerequisites

- Go 1.18+ (for building from source)
- Kagi API key (available to Kagi subscribers)

## Installation

### Using Go

```bash
# Clone the repository
git clone https://github.com/mgomes/kagimcp.git
cd kagimcp

# Install dependencies
go mod download

# Build the application
go build -o kagimcp
```

## Usage

### Command-Line Options

```
Usage of ./kagi-mcp:
  -api-key string
        Kagi API key (can also be set with KAGI_API_KEY environment variable)
  -port string
        Port for SSE server (default "8080")
  -t string
        Transport type (stdio or sse) (default "stdio")
```

### Running in Stdio Mode

This mode is useful for direct integration with LLM platforms that support subprocess communication.

```bash
# Using direct binary
KAGI_API_KEY=your_api_key ./kagimcp -t stdio

### Running in SSE Mode

This mode starts an HTTP server that communicates using Server-Sent Events (SSE).

```bash
# Using direct binary
KAGI_API_KEY=your_api_key ./kagimcp -t sse -port 8080

## Available Tools

### Kagi Search

Searches the web using Kagi Search API.

Parameters:
- `query` (string, required): The search query string
- `limit` (number, optional): Maximum number of results (1-10, default: 5)

Example:
```json
{
  "name": "kagi_search",
  "arguments": {
    "query": "climate change solutions",
    "limit": 3,
    "type": "news"
  }
}
```

### Kagi Summarize

Summarizes a webpage using Kagi's FastGPT API.

Parameters:
- `url` (string, required): URL of the webpage to summarize
- `engine` (string, optional): Summarization engine to use ("cecil", "agnes", or "muriel" default: "agnes"),
- `summary_type` (string, optional): Type of summary to generate ("summary", "takeaway" default: "summary")

Example:
```json
{
  "name": "kagi_summarize",
  "arguments": {
    "url": "https://en.wikipedia.org/wiki/Artificial_intelligence",
    "engine": "cecil",
    "summary_type": "summary"
  }
}
```

## Integration Examples

### Integrating with Claude

You can connect Claude to this MCP server to give it the ability to search the web and summarize web pages.

1. Start the server in stdio mode
2. Configure Claude to use it as an MCP tool

```json
{
  "mcpServers": {
    "kagi": {
      "command": "./kagimcp",
      "args": [],
      "env": {
        "KAGI_API_KEY": "YOUR_API_KEY_HERE"
      }
    }
  }
}
```

### Integrating with Other LLM Platforms

This server is compatible with any LLM platform that supports the Model Context Protocol. Refer to your platform's documentation for specific integration steps.

## License

MIT
