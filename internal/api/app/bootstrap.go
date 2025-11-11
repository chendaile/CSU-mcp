package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astaxie/beego"
)

const (
	configPath = "configs/api/conf/app.conf"
	viewsPath  = "web/views"
	staticPath = "web/static"
	logPath    = "var/logs/api/project.log"
)

// Initialize wires configuration, static assets, and logging so that
// Beego can serve the API regardless of the current working directory.
func Initialize() error {
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}

	beego.BConfig.WebConfig.ViewsPath = viewsPath
	if err := beego.AddViewPath(viewsPath); err != nil {
		return fmt.Errorf("register views: %w", err)
	}
	beego.SetStaticPath("/static", staticPath)

	if err := beego.LoadAppConfig("ini", configPath); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	loggerCfg := fmt.Sprintf(`{"filename":"%s","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`, logPath)
	if err := beego.SetLogger("file", loggerCfg); err != nil {
		return fmt.Errorf("configure logger: %w", err)
	}

	return nil
}
