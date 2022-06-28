// Package types provides publicly available type for the weather package
package types

import (
	"sync"
	"time"
)

type Coordinates struct {
	Lon float64
	Lat float64
}

type PrecipitationRecord struct {
	PrecipitationType string
	Window            int32
	Amount            int32
}

type WindData struct {
	Speed float64
	Deg   int32
	Gust  float64
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

type Config struct {
	ApiKey          string
	DefaultLocation Coordinates
	Service         string
}

type GetWeatherFunc func(Config, chan CurrentWeather, *sync.WaitGroup)
