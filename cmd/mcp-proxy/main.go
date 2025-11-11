package main

import (
	"net/http"
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"campusapp/internal/mcp/campus"
	"campusapp/internal/mcp/config"
	"campusapp/internal/mcp/httpserver"
	"campusapp/internal/mcp/tools"
)

func main() {
	initLogger()
	cfg, err := config.Load()
	if err != nil {
		logs.Critical("load config: %v", err)
		os.Exit(1)
	}

	client, err := campus.NewClient(cfg.BaseURL, cfg.Token, cfg.Timeout)
	if err != nil {
		logs.Critical("init campus client: %v", err)
		os.Exit(1)
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

	logs.Info("csu MCP server listening on %s (csugo base %s)", cfg.HTTPAddr, cfg.BaseURL)

	appHandler := httpserver.LandingMiddleware(handler, cfg.ImplementationName, cfg.BaseURL, cfg.HTTPAddr)

	if err := http.ListenAndServe(cfg.HTTPAddr, httpserver.LoggingMiddleware(appHandler)); err != nil {
		logs.Critical("http server stopped: %v", err)
		os.Exit(1)
	}
}

func initLogger() {
	logs.SetLogger(logs.AdapterConsole)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(2)
	logs.SetLevel(logs.LevelDebug)
}
