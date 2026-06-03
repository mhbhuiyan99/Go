package models

import (
	"math/rand"
	"math"
)

type Weather struct {
	City string `json:"city"`
	Temperature float64 `json:"temperature"`
}

func GetWeather(city string) Weather {
	return Weather{
		City: city,
		Temperature: math.Round(rand.Float64() * 38 *100) / 100, // Random temperature between 0 and 38
	}
}