// Package weather acts as a single point of contact for acquiring
// weather data from any available source
package weather

import (
	"shearwater.ai/pocs/weather/internal/adapters/open_weather"
	"shearwater.ai/pocs/weather/types"
	"sync"
)

func GetWeather(config types.Config, respCh chan types.CurrentWeather, wg *sync.WaitGroup) {
	var weatherFunc types.GetWeatherFunc
	switch config.Service {
	case "OpenWeather":
		weatherFunc = weather.GetOpenWeather
	}
	if weatherFunc == nil {
		close(respCh)
		wg.Done()
	} else {
		weatherFunc(config, respCh, wg)
	}
}
