package controllers

import (
	"net/http"
	"strings"

	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

type Error struct {
	StateCode int
	Error     string
}

func (c *ErrorController) Error404() {
	payload := Error{
		StateCode: 404,
		Error:     "api not found",
	}
	c.Data["json"] = payload

	path := c.Ctx.Input.URL()
	accept := c.Ctx.Request.Header.Get("Accept")
	wantsJSON := strings.HasPrefix(path, "/api/") ||
		strings.HasPrefix(path, "/wxapp") ||
		strings.Contains(accept, "application/json")

	if wantsJSON {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.ServeJSON()
		c.StopRun()
		return
	}

	c.TplName = "errors/404.html"
}
