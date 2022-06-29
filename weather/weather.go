// Package weather acts as a single point of contact for acquiring
// weather data from any available source
package weather

import (
	"github.com/jpgringo/shearwater_pocs_external_sources/weather/internal/adapters/open_weather"
	"github.com/jpgringo/shearwater_pocs_external_sources/weather/type_definitions"
	"sync"
)

func GetWeather(config type_definitions.Config, respCh chan type_definitions.CurrentWeather, wg *sync.WaitGroup) {
	var weatherFunc type_definitions.GetWeatherFunc
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
