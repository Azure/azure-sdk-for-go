package main

import (
	"fmt"
	"maps/weather"
)

func main() {
	fmt.Print("Hello Go")
	weather.NewWeatherClient(nil, nil)
}
