package weather

import (
	"log"
	"sync"
	"testing"
)

func TestGetOpenWeatherUpdate(t *testing.T) {
	GetOpenWeatherUpdate()
}

func TestGetWeather(t *testing.T) {
	wg := sync.WaitGroup{}
	respChan := make(chan CurrentWeather)
	wg.Add(1)
	go GetWeather(respChan, &wg)
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
	go GetWeather(respChan, &wg)
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
