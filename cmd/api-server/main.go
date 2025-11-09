package main

import (
	"log"

	"github.com/astaxie/beego"

	"campusapp/internal/api/app"
	"campusapp/internal/api/controllers"
	_ "campusapp/internal/api/routers"
)

func main() {
	if err := app.Initialize(); err != nil {
		log.Fatal(err)
	}

	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()
}
