package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const owApiKey string = "0312e9e8566e6d2591dabc1e779c60a7"
const defaultLat = 43.64913651147442
const defaultLon = -79.45198018043132
const owRequestTemplate = "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric" // &exclude={part}

var lastFetched *CurrentWeather
var weatherLock sync.Mutex

type WindData struct {
	Speed float64
	Deg   int32
	Gust  float64
}

type Coordinates struct {
	Lon float64
	Lat float64
}

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
	Loc        Coordinates `json:"coord"`
	Name       string
	Main       owMain
	Visibility int32
	Wind       WindData
	Clouds     owClouds
	Rain       owPrecip
	Snow       owPrecip
	Dt         int64
}

type CurrentWeather struct {
	Loc           Coordinates
	Name          string
	Timestamp     time.Time
	Temp          float64
	Humidity      int32
	Pressure      int32
	Visibility    int32
	Wind          WindData
	Clouds        int32
	Precipitation []PrecipitationRecord
}

type PrecipitationRecord struct {
	PrecipitationType string
	Window            int32
	Amount            int32
}

func (v openWeatherData) convertToCurrentWeather() CurrentWeather {
	cw := CurrentWeather{
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
		PrecipitationRecord{
			PrecipitationType: "Rain",
			Window:            1,
			Amount:            v.Rain.Last1Hour,
		},
		PrecipitationRecord{
			PrecipitationType: "Rain",
			Window:            3,
			Amount:            v.Rain.Last3Hours,
		},
		PrecipitationRecord{
			PrecipitationType: "Snow",
			Window:            1,
			Amount:            v.Snow.Last1Hour,
		},
		PrecipitationRecord{
			PrecipitationType: "Snow",
			Window:            3,
			Amount:            v.Snow.Last3Hours,
		},
	)
	return cw
}

func GetOpenWeatherUpdate() (cw CurrentWeather) {
	req := fmt.Sprintf(owRequestTemplate, defaultLat, defaultLon, owApiKey)
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

func GetWeather(respCh chan CurrentWeather, wg *sync.WaitGroup) {
	weatherLock.Lock()
	defer weatherLock.Unlock()
	if lastFetched == nil {
		log.Println("weather has not yet been fetched")
		cw := GetOpenWeatherUpdate()
		lastFetched = &cw
	} else {
		log.Println("already got the weather!!")
	}
	respCh <- *lastFetched
	defer wg.Done()
}
