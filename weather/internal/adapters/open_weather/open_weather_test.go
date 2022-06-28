package weather

import (
	"log"
	"shearwater.ai/pocs/weather/types"
	"sync"
	"testing"
)

func defaultConfig() types.Config {
	return types.Config{
		ApiKey: "0312e9e8566e6d2591dabc1e779c60a7",
		DefaultLocation: types.Coordinates{
			Lat: 43.64913651147442,
			Lon: -79.45198018043132,
		},
	}
}

func TestGetOpenWeatherUpdate(t *testing.T) {
	getOpenWeatherUpdate(defaultConfig())
}

func TestGetWeather(t *testing.T) {
	wg := sync.WaitGroup{}
	respChan := make(chan types.CurrentWeather)
	wg.Add(1)
	config := defaultConfig()
	go GetOpenWeather(config, respChan, &wg)
receiveLoop:
	for {
		select {
		case cw := <-respChan:
			log.Println("received weather:", cw)
			break receiveLoop
		}
	}
	wg.Wait()
	wg.Add(1)
	go GetOpenWeather(config, respChan, &wg)
loopTwo:
	for {
		select {
		case cw := <-respChan:
			log.Println("received weather:", cw)
			break loopTwo
		}
	}
	wg.Wait()
}
