package main

import (
	"fmt"
	"net/http"
	"time"
)

func landingMiddleware(next http.Handler, implName, baseURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>%s</title>
  <style>
    body { font-family: sans-serif; background: #0f172a; color: #f8fafc; display: flex; align-items: center; justify-content: center; min-height: 100vh; margin: 0; }
    .card { background: #1e293b; padding: 32px; border-radius: 12px; box-shadow: 0 20px 40px rgba(0,0,0,0.35); max-width: 520px; }
    h1 { margin-top: 0; font-size: 1.6rem; }
    p { margin: 0 0 12px; line-height: 1.5; }
    small { color: #94a3b8; }
    code { background: #0f172a; padding: 2px 6px; border-radius: 4px; }
  </style>
</head>
<body>
  <div class="card">
    <h1>%s 已就绪</h1>
    <p>模型上下文协议 (MCP) 代理正在运行，并已连接到 <code>%s</code>。</p>
    <p>如果你能看到这个页面，说明成功访问到了 MCP Server 的 HTTP 入口。</p>
    <small>更新时间：%s</small>
  </div>
</body>
</html>`, implName, implName, baseURL, time.Now().Format(time.RFC3339))
			return
		}
		next.ServeHTTP(w, r)
	})
}
