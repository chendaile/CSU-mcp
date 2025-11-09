package main

import (
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"campusapp/internal/mcp/campus"
	"campusapp/internal/mcp/config"
	"campusapp/internal/mcp/tools"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	client, err := campus.NewClient(cfg.BaseURL, cfg.Token, cfg.Timeout)
	if err != nil {
		log.Fatal(err)
	}

	toolset := &tools.Toolset{Client: client}

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		server := mcp.NewServer(&mcp.Implementation{
			Name:    cfg.ImplementationName,
			Version: cfg.ImplementationVersion,
		}, nil)
		toolset.Register(server)
		return server
	}, nil)

	log.Printf("csu MCP server listening on %s (csugo base %s)", cfg.HTTPAddr, cfg.BaseURL)

	appHandler := landingMiddleware(handler, cfg.ImplementationName, cfg.BaseURL)

	if err := http.ListenAndServe(cfg.HTTPAddr, loggingMiddleware(appHandler)); err != nil {
		log.Fatal(err)
	}
}
