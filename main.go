package main

import (
	"fmt"
	"log"
	"os"

	forecast "github.com/mlbright/forecast/v2"
)

func main() {
	key := os.Getenv("FORECAST_APIKEY")
	lat := os.Args[1]
	long := os.Args[2]

	f, err := forecast.Get(key, lat, long, "now", forecast.US)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %f", f.Minutely.Data[0].Summary, f.Minutely.Data[0].PrecipProbability)

}
