package weather

import (
	"encoding/json"
	"fmt"
	"github.com/jpgringo/shearwater_pocs_external_sources/weather/type_definitions"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const owRequestTemplate = "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric" // &exclude={part}

var lastFetched *type_definitions.CurrentWeather
var weatherLock sync.Mutex

type owMain struct {
	Temp     float64
	Humidity int32
	Pressure int32
}

type owPrecipitation struct {
	Last1Hour  float64 `json:"1h"`
	Last3Hours float64 `json:"3h"`
}

type owClouds struct {
	All int32
}

type openWeatherData struct {
	Base       string
	Loc        type_definitions.Coordinates `json:"coord"`
	Name       string
	Main       owMain
	Visibility int32
	Wind       type_definitions.WindData
	Clouds     owClouds
	Rain       owPrecipitation
	Snow       owPrecipitation
	Dt         int64
}

func (v openWeatherData) convertToCurrentWeather() type_definitions.CurrentWeather {
	cw := type_definitions.CurrentWeather{
		Loc:        v.Loc,
		Name:       v.Name,
		Timestamp:  time.Unix(v.Dt, 0),
		Temp:       v.Main.Temp,
		Humidity:   v.Main.Humidity,
		Pressure:   v.Main.Pressure,
		Visibility: v.Visibility,
		Wind:       v.Wind,
		Clouds:     v.Clouds.All,
	}
	cw.Precipitation = append(cw.Precipitation,
		type_definitions.PrecipitationRecord{
			PrecipitationType: "Rain",
			Window:            1,
			Amount:            v.Rain.Last1Hour,
		},
		type_definitions.PrecipitationRecord{
			PrecipitationType: "Rain",
			Window:            3,
			Amount:            v.Rain.Last3Hours,
		},
		type_definitions.PrecipitationRecord{
			PrecipitationType: "Snow",
			Window:            1,
			Amount:            v.Snow.Last1Hour,
		},
		type_definitions.PrecipitationRecord{
			PrecipitationType: "Snow",
			Window:            3,
			Amount:            v.Snow.Last3Hours,
		},
	)
	return cw
}

func getOpenWeatherUpdate(config type_definitions.Config) (cw type_definitions.CurrentWeather) {
	req := fmt.Sprintf(owRequestTemplate, config.DefaultLocation.Lat, config.DefaultLocation.Lon, config.ApiKey)
	resp, err := http.Get(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var ow openWeatherData
	err = json.Unmarshal(body, &ow)
	if err != nil {
		log.Fatal(err)
	}
	cw = ow.convertToCurrentWeather()
	return cw
}

func GetOpenWeather(config type_definitions.Config, respCh chan type_definitions.CurrentWeather, wg *sync.WaitGroup) {
	weatherLock.Lock()
	defer weatherLock.Unlock()
	if lastFetched == nil {
		log.Println("weather has not yet been fetched")
		cw := getOpenWeatherUpdate(config)
		lastFetched = &cw
	} else {
		log.Println("already got the weather!!")
	}
	respCh <- *lastFetched
	defer wg.Done()
}
