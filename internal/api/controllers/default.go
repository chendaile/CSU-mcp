package controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego"

	"campusapp/internal/branding"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	appName := beego.AppConfig.String("appname")
	if appName == "" {
		appName = "CSU Campus API"
	}
	httpPort := beego.AppConfig.String("httpport")
	if httpPort == "" {
		httpPort = fmt.Sprintf("%d", beego.BConfig.Listen.HTTPPort)
	}

	scheme := c.Ctx.Input.Scheme()
	if scheme == "" {
		scheme = "http"
	}
	hostHeader := c.Ctx.Input.Host()
	hostName, hostPort := branding.ResolveHost(hostHeader, "127.0.0.1", httpPort)
	displayHost := branding.JoinHostPort(hostName, hostPort)
	metaAddr := fmt.Sprintf("%s://%s", scheme, displayHost)

	mcpPort := branding.PortFromAddress(os.Getenv("MCP_HTTP_ADDR"), "13000")
	mcpURL := fmt.Sprintf("%s://%s", scheme, branding.JoinHostPort(hostName, mcpPort))

	c.Data["Title"] = appName
	c.Data["Subtitle"] = "REST & 数据聚合服务"
	c.Data["Description"] = "基于 Beego 的校园数据聚合 API，提供成绩、课表、图书、自习室等统一接口。"
	c.Data["Badge"] = "API 服务"
	c.Data["Meta"] = fmt.Sprintf("监听地址：%s", metaAddr)
	c.Data["Highlights"] = []string{
		"GET /api/v1/jwc/:id/:pwd/grade — 查询成绩",
		"GET /api/v1/jwc/:id/:pwd/rank — 查询综合排名",
		"GET /api/v1/jwc/:id/:pwd/class/:term/:week — 获取课表",
		"GET /api/v1/bus/search/:start/:end/:day — 校车查询",
	}
	c.Data["PrimaryAction"] = "查看 GitHub"
	c.Data["PrimaryURL"] = branding.AuthorURL
	c.Data["SecondaryAction"] = "访问 MCP 代理"
	c.Data["SecondaryURL"] = mcpURL
	c.Data["AuthorName"] = branding.AuthorName
	c.Data["AuthorURL"] = branding.AuthorURL
	c.Data["Timestamp"] = time.Now().Format("2006-01-02 15:04:05 MST")
	c.TplName = "index.html"
}
