package main

import (
	_ "weather-dashboard/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

