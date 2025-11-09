package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"

	"campusapp/internal/branding"
)

type landingPageData struct {
	Title           string
	Subtitle        string
	Description     string
	Badge           string
	Highlights      []string
	PrimaryAction   string
	PrimaryURL      string
	SecondaryAction string
	SecondaryURL    string
	AuthorName      string
	AuthorURL       string
	Meta            string
	Timestamp       string
}

var (
	landingTemplateOnce sync.Once
	landingTemplate     *template.Template
	landingTemplateErr  error
)

func landingMiddleware(next http.Handler, implName, baseURL, listenAddr string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/" {
			scheme := resolveScheme(r)
			listenPort := branding.PortFromAddress(listenAddr, "13000")
			hostName, hostPort := branding.ResolveHost(r.Host, "127.0.0.1", listenPort)
			displayHost := branding.JoinHostPort(hostName, hostPort)
			metaAddr := fmt.Sprintf("%s://%s", scheme, displayHost)
			apiURL := branding.RewriteBaseURL(baseURL, scheme, hostName)

			data := landingPageData{
				Title:       fmt.Sprintf("%s 已就绪", implName),
				Subtitle:    "Model Context Protocol 代理",
				Description: "MCP HTTP 入口正在运行，可通过客户端连接并调用 CSUGO 工具集。",
				Badge:       "MCP 服务",
				Highlights: []string{
					fmt.Sprintf("MCP 监听：%s", metaAddr),
					fmt.Sprintf("目标 API：%s", apiURL),
					"可用工具：csu.grade · csu.rank · csu.classes · csu.bus_search · csu.jobs",
				},
				PrimaryAction:   "查看 MCP 说明",
				PrimaryURL:      "https://modelcontextprotocol.io",
				SecondaryAction: "访问 API 服务",
				SecondaryURL:    apiURL,
				AuthorName:      branding.AuthorName,
				AuthorURL:       branding.AuthorURL,
				Meta:            fmt.Sprintf("监听地址：%s", metaAddr),
				Timestamp:       time.Now().Format("2006-01-02 15:04:05 MST"),
			}
			if err := renderLanding(w, data); err != nil {
				logs.Error("render landing: %v", err)
				http.Error(w, "landing page unavailable", http.StatusInternalServerError)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func renderLanding(w http.ResponseWriter, data landingPageData) error {
	tmpl, err := loadLandingTemplate()
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "index.html", data)
}

func loadLandingTemplate() (*template.Template, error) {
	landingTemplateOnce.Do(func() {
		landingTemplate, landingTemplateErr = template.ParseFiles(branding.LandingTemplatePath)
	})
	return landingTemplate, landingTemplateErr
}

func resolveScheme(r *http.Request) string {
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
