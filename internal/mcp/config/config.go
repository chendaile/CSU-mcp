package config

import (
	"fmt"
	"os"
	"strings"
	"time"
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
	cfg := Server{
		BaseURL:               getEnv("CSUGO_BASE_URL", "http://localhost:12000"),
		Token:                 getEnv("CSUGO_TOKEN", "csugo-token"),
		Timeout:               15 * time.Second,
		ImplementationName:    getEnv("MCP_IMPLEMENTATION_NAME", "csu-mcp-proxy"),
		ImplementationVersion: getEnv("MCP_IMPLEMENTATION_VERSION", "0.1.0"),
		HTTPAddr:              getEnv("MCP_HTTP_ADDR", ""),
	}

	if cfg.BaseURL == "" {
		return cfg, fmt.Errorf("CSUGO_BASE_URL must be set")
	}

	if strings.TrimSpace(cfg.HTTPAddr) == "" {
		cfg.HTTPAddr = ":13000"
	}

	if raw := os.Getenv("CSUGO_TIMEOUT"); raw != "" {
		timeout, err := time.ParseDuration(raw)
		if err != nil {
			return cfg, fmt.Errorf("invalid CSUGO_TIMEOUT: %w", err)
		}
		if timeout <= 0 {
			return cfg, fmt.Errorf("CSUGO_TIMEOUT must be > 0")
		}
		cfg.Timeout = timeout
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
