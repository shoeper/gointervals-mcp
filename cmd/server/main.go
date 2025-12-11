package main

import (
	"fmt"
	"github.com/shoeper/gointervals-mcp/internal/config"
	"github.com/shoeper/gointervals-mcp/internal/tools"
	"net/http"

	"github.com/mark3labs/mcp-go/server"
)

type Activity struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	StartDateLocal string  `json:"start_date_local"`
	Type           string  `json:"type"`
	MovingTime     int     `json:"moving_time"`
	Distance       float64 `json:"distance"`
	TotalElevation float64 `json:"total_elevation_gain"`
}

// --- Auth Middleware ---
// This middleware now wraps the SINGLE Streamable HTTP endpoint
func authMiddleware(next http.Handler, validToken string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+validToken {
			http.Error(w, "Unauthorized: Invalid Bearer Token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// --- Main Entry Point ---
func main() {
	cfg := config.Load()

	// 1. Create Base Server
	mcpServer := server.NewMCPServer("gointervals-mcp", "1.0.0")
    tools.RegisterAll(mcpServer, &cfg)

	// 2. Create STREAMABLE HTTP Server, implements the single-endpoint specification.
	streamableServer := server.NewStreamableHTTPServer(mcpServer)

	// 3. Setup Router
	mux := http.NewServeMux()
	
	// We expose a SINGLE endpoint "/mcp" that handles both POST (messages) and GET (connect)
	mux.Handle("/mcp", authMiddleware(streamableServer, cfg.McpAuthToken))

	// Optional: Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Printf("Streamable HTTP MCP Server listening on 0.0.0.0:%s/mcp\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		panic(err)
	}
}