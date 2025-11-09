package routers

import (
	"campusapp/internal/api/controllers"
	_ "campusapp/internal/api/middleware"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1",
			beego.NSInclude(&controllers.JwcController{}),
			beego.NSInclude(&controllers.BusController{}),
			beego.NSInclude(&controllers.JobController{}),
			beego.NSInclude(&controllers.CetController{}),
			beego.NSInclude(&controllers.LibController{}),
			beego.NSInclude(&controllers.ClassRoomController{}),
			beego.NSInclude(&controllers.NewsController{}),
		),
	)
	wx := beego.NewNamespace("/wxapp",
		beego.NSInclude(&controllers.WxUserController{}),
	)
	beego.AddNamespace(wx)
	beego.AddNamespace(ns)
}
