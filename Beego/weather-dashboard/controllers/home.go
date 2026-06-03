package controllers

import (
	"sync"
	"weather-dashboard/models"
	beego "github.com/beego/beego/v2/server/web"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	cities := []string{
		"Dhaka",
		"Chittagong",
		"Khulna",
		"Rajshahi",
		"Sylhet",
	}

	var wg sync.WaitGroup

	weatherList := make([]models.Weather, len(cities))

	for i, city := range cities {
		wg.Add(1)

		go func(i int, city string) {
			defer wg.Done()

			weatherList[i] = models.GetWeather(city)
		} (i, city)
	}

	wg.Wait()

	c.Data["Weather"] = weatherList
	c.TplName = "index.tpl"
}