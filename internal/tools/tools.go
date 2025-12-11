package tools

import (
	"intervalsicu-mcp/internal/config"

	"github.com/mark3labs/mcp-go/server"
)

// RegisterAll registers all tools with the MCP server
func RegisterAll(s *server.MCPServer, config *config.Config) {
    GetActivities(s, config)
}
