package routers

import (
	"weather-dashboard/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.HomeController{})
}
