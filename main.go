package main

import (
	"flag"
	"log/slog"
	"os"
)

func main() {
	var transport string
	var apiKey string
	var port string

	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&apiKey, "api-key", "", "Kagi API key")
	flag.StringVar(&port, "port", "8080", "Port for SSE server")
	flag.Parse()

	if apiKey == "" {
		apiKey = os.Getenv("KAGI_API_KEY")
		if apiKey == "" {
			slog.Error("Kagi API key not provided. Use -api-key flag or KAGI_API_KEY environment variable")
			os.Exit(1)
		}
	}

	kagiSrv := NewKagiServer(apiKey)

	switch transport {
	case "stdio":
		slog.Info("Starting stdio server")
		if err := kagiSrv.ServeStdio(); err != nil {
			slog.Error("Server error", "err", err)
			os.Exit(1)
		}
	case "sse":
		addr := ":" + port
		slog.Info("Starting SSE server", "addr", addr)
		if err := kagiSrv.ServeSSE(addr); err != nil {
			slog.Error("Server error", "err", err)
			os.Exit(1)
		}
	default:
		slog.Error("Unknown transport. Use 'stdio' or 'sse'", "transport", transport)
		os.Exit(1)
	}
}
