// Package mcp provides MCP (Model Context Protocol) server for prism-maturity.
package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server provides MCP tools for maturity management.
type Server struct {
	dataDir string
}

// NewServer creates a new MCP server.
func NewServer() *Server {
	return &Server{}
}

// WithDataDir sets the data directory for the server.
func (s *Server) WithDataDir(dir string) *Server {
	s.dataDir = dir
	return s
}

// RegisterTools registers all maturity tools with the MCP server.
func (s *Server) RegisterTools(server *mcp.Server) {
	s.RegisterCategoryTools(server)
	s.RegisterGoalTools(server)
	s.RegisterMaturityTools(server)
}

// Run starts the MCP server.
func (s *Server) Run(ctx context.Context) error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "prism-maturity",
		Version: "0.1.0",
	}, nil)

	s.RegisterTools(server)

	return server.Run(ctx, &mcp.StdioTransport{})
}
