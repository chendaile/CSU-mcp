package branding

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// ResolveHost splits a raw host header into host & port, falling back if absent.
func ResolveHost(raw, fallbackHost, fallbackPort string) (string, string) {
	host := fallbackHost
	port := fallbackPort

	if strings.TrimSpace(raw) == "" {
		return host, port
	}

	if h, p, err := net.SplitHostPort(raw); err == nil {
		if h != "" {
			host = h
		}
		if p != "" {
			port = p
		}
		return host, port
	}

	// raw doesn't contain a port
	return raw, port
}

// JoinHostPort builds host[:port] while respecting IPv6 formatting.
func JoinHostPort(host, port string) string {
	if host == "" {
		host = "127.0.0.1"
	}
	port = strings.TrimSpace(port)
	if port == "" {
		return host
	}
	return net.JoinHostPort(host, port)
}

// PortFromAddress extracts the port from an address or returns fallback.
func PortFromAddress(addr, fallback string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return fallback
	}
	if strings.HasPrefix(addr, ":") && len(addr) > 1 {
		return addr[1:]
	}
	if strings.Count(addr, ":") == 0 && strings.Index(addr, ".") == -1 {
		// looks like just a port number
		return addr
	}
	if _, p, err := net.SplitHostPort(addr); err == nil && p != "" {
		return p
	}
	return fallback
}

// RewriteBaseURL takes an existing base URL and rewrites its host with the provided host.
func RewriteBaseURL(base, scheme, host string) string {
	if strings.TrimSpace(host) == "" {
		host = "127.0.0.1"
	}
	parsed, err := url.Parse(base)
	if err != nil {
		if scheme == "" {
			scheme = "http"
		}
		return fmt.Sprintf("%s://%s", scheme, host)
	}

	if scheme != "" {
		parsed.Scheme = scheme
	} else if parsed.Scheme == "" {
		parsed.Scheme = "http"
	}

	port := parsed.Port()
	if port == "" {
		if _, p, err := net.SplitHostPort(parsed.Host); err == nil {
			port = p
		}
	}

	parsed.Host = JoinHostPort(host, port)
	return parsed.String()
}
