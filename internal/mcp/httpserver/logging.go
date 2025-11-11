package httpserver

import (
	"net"
	"net/http"
	"time"

	"github.com/astaxie/beego/logs"
)

// LoggingMiddleware captures status and latency for every MCP HTTP request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		host := r.RemoteAddr
		if parsedHost, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			host = parsedHost
		}

		logs.Info("| %15s | %3d | %12s | %s %s", host, lrw.status, duration, r.Method, r.URL.Path)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (lrw *loggingResponseWriter) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.status == 0 {
		lrw.status = http.StatusOK
	}
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytes += n
	return n, err
}
