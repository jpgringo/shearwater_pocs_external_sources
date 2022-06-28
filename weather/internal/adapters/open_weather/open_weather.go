package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"shearwater.ai/pocs/weather/types"
	"sync"
	"time"
)

const owApiKey string = "0312e9e8566e6d2591dabc1e779c60a7"
const defaultLat = 43.64913651147442
const defaultLon = -79.45198018043132
const owRequestTemplate = "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric" // &exclude={part}

var lastFetched *types.CurrentWeather
var weatherLock sync.Mutex

type owMain struct {
	Temp     float64
	Humidity int32
	Pressure int32
}

type owPrecip struct {
	Last1Hour  int32 `json:"1h"`
	Last3Hours int32 `json:"3h"`
}

type owClouds struct {
	All int32
}

type openWeatherData struct {
	Base       string
	Loc        types.Coordinates `json:"coord"`
	Name       string
	Main       owMain
	Visibility int32
	Wind       types.WindData
	Clouds     owClouds
	Rain       owPrecip
	Snow       owPrecip
	Dt         int64
}

func (v openWeatherData) convertToCurrentWeather() types.CurrentWeather {
	cw := types.CurrentWeather{
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
		types.PrecipitationRecord{
			PrecipitationType: "Rain",
			Window:            1,
			Amount:            v.Rain.Last1Hour,
		},
		types.PrecipitationRecord{
			PrecipitationType: "Rain",
			Window:            3,
			Amount:            v.Rain.Last3Hours,
		},
		types.PrecipitationRecord{
			PrecipitationType: "Snow",
			Window:            1,
			Amount:            v.Snow.Last1Hour,
		},
		types.PrecipitationRecord{
			PrecipitationType: "Snow",
			Window:            3,
			Amount:            v.Snow.Last3Hours,
		},
	)
	return cw
}

func getOpenWeatherUpdate(config types.Config) (cw types.CurrentWeather) {
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

func GetOpenWeather(config types.Config, respCh chan types.CurrentWeather, wg *sync.WaitGroup) {
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
