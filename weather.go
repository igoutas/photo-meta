package main

import (
	"encoding/json"
	"fmt"
)

type WeatherClient struct {
	apiKey     string
	httpClient HttpClientInterface
	cache      map[string]*WeatherMeta
}

const (
	baseUrl string = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/"
)

func (weather *WeatherClient) GetWeatherCondition(latitude, longitude, datetime string, c chan *WeatherMeta) {
	cacheKey := latitude + longitude + datetime
	res, ok := weather.cache[cacheKey]
	if ok {
		c <- res
	} else {
		var weatherMeta *WeatherMeta
		url := baseUrl + latitude + "," + longitude + "/" + datetime + "?key=" + weather.apiKey + "&include=current&unitGroup=uk"

		body, err := weather.httpClient.Get(url)
		if err == nil {
			var weatherData WeatherData
			err = json.Unmarshal(body, &weatherData)

			if err == nil {
				weatherMeta = &WeatherMeta{
					WeatherData: &weatherData,
					Error:       nil,
				}
			}
		}

		if err != nil {
			fmt.Println("Error while calling the weather api")
			weatherMeta = &WeatherMeta{
				WeatherData: nil,
				Error:       err,
			}
		}

		// update the cache
		weather.cache[cacheKey] = weatherMeta
		c <- weatherMeta
	}
}

func NewWeatherClient(httpClient HttpClientInterface, apiKey string) WeatherClientInterface {
	weatherClient := WeatherClient{}
	weatherClient.apiKey = apiKey
	weatherClient.httpClient = httpClient
	weatherClient.cache = make(map[string]*WeatherMeta)

	return &weatherClient
}
