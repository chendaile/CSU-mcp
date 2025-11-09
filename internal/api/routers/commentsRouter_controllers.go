package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["campusapp/internal/api/controllers:BusController"] = append(beego.GlobalControllerRouter["campusapp/internal/api/controllers:BusController"],
		beego.ControllerComments{
			Method:           "Search",
			Router:           `/bus/search/:start/:end/:time`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["campusapp/internal/api/controllers:CetController"] = append(beego.GlobalControllerRouter["campusapp/internal/api/controllers:CetController"],
		beego.ControllerComments{
			Method:           "GetHGrade",
			Router:           `/cet/hgrade/:id/:name`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["campusapp/internal/api/controllers:CetController"] = append(beego.GlobalControllerRouter["campusapp/internal/api/controllers:CetController"],
		beego.ControllerComments{
			Method:           "GetZKZ",
			Router:           `/cet/zkz/:id/:type`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["campusapp/internal/api/controllers:JobController"] = append(beego.GlobalControllerRouter["campusapp/internal/api/controllers:JobController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/job/:typeid/:pageindex/:pagesize/:hastime`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["campusapp/internal/api/controllers:JwcController"] = append(beego.GlobalControllerRouter["campusapp/internal/api/controllers:JwcController"],
		beego.ControllerComments{
			Method:           "Class",
			Router:           `/jwc/:id/:pwd/class/:term/:week`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["campusapp/internal/api/controllers:JwcController"] = append(beego.GlobalControllerRouter["campusapp/internal/api/controllers:JwcController"],
		beego.ControllerComments{
			Method:           "Grade",
			Router:           `/jwc/:id/:pwd/grade`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["campusapp/internal/api/controllers:JwcController"] = append(beego.GlobalControllerRouter["campusapp/internal/api/controllers:JwcController"],
		beego.ControllerComments{
			Method:           "Rank",
			Router:           `/jwc/:id/:pwd/rank`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["campusapp/internal/api/controllers:LibController"] = append(beego.GlobalControllerRouter["campusapp/internal/api/controllers:LibController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           `/lib/list/:id/:pwd`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["campusapp/internal/api/controllers:LibController"] = append(beego.GlobalControllerRouter["campusapp/internal/api/controllers:LibController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `/lib/login/:id/:pwd`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["campusapp/internal/api/controllers:WxUserController"] = append(beego.GlobalControllerRouter["campusapp/internal/api/controllers:WxUserController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

}
