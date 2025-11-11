package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/astaxie/beego/config"
)

type Server struct {
	BaseURL               string
	Token                 string
	Timeout               time.Duration
	ImplementationName    string
	ImplementationVersion string
	HTTPAddr              string
}

func Load() (Server, error) {
	cfgPath, err := resolveConfigPath()
	if err != nil {
		return Server{}, err
	}

	rawCfg, err := config.NewConfig("ini", cfgPath)
	if err != nil {
		return Server{}, fmt.Errorf("load mcp config %s: %w", cfgPath, err)
	}

	timeoutStr := rawCfg.DefaultString("MCP::Timeout", "15s")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return Server{}, fmt.Errorf("invalid MCP::Timeout %q: %w", timeoutStr, err)
	}
	if timeout <= 0 {
		return Server{}, fmt.Errorf("MCP::Timeout must be > 0")
	}

	cfg := Server{
		BaseURL:               strings.TrimSpace(rawCfg.DefaultString("MCP::BaseURL", "")),
		Token:                 rawCfg.DefaultString("MCP::Token", ""),
		Timeout:               timeout,
		ImplementationName:    rawCfg.DefaultString("MCP::ImplementationName", "csu-mcp-proxy"),
		ImplementationVersion: rawCfg.DefaultString("MCP::ImplementationVersion", "0.1.0"),
		HTTPAddr:              strings.TrimSpace(rawCfg.DefaultString("MCP::HTTPAddr", ":13000")),
	}

	if cfg.BaseURL == "" {
		return cfg, fmt.Errorf("MCP::BaseURL must be set in %s", cfgPath)
	}
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":13000"
	}

	return cfg, nil
}

func resolveConfigPath() (string, error) {
	start, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("determine working directory: %w", err)
	}
	target := filepath.Join("configs", "api", "conf", "app.conf")
	dir := start
	for {
		candidate := filepath.Join(dir, target)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	alt := filepath.Join("/app", target)
	if _, err := os.Stat(alt); err == nil {
		return alt, nil
	}
	return "", fmt.Errorf("unable to locate %s from %s", target, start)
}
