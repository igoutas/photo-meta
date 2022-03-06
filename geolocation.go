package main

import (
	"encoding/json"
	"fmt"
)

type GeolocationResolver struct {
	apiKey     string
	httpClient HttpClientInterface
	cache      map[string]*GeoLocationMeta
}

func (geo *GeolocationResolver) GetGeoLocationMeta(latitude, longitude string, c chan *GeoLocationMeta) {
	cacheKey := latitude + longitude
	res, ok := geo.cache[cacheKey]
	if ok {
		c <- res
	} else {
		var geoMeta *GeoLocationMeta
		url := "https://revgeocode.search.hereapi.com/v1/revgeocode?apiKey=" + geo.apiKey + "&at=" + latitude + "," + longitude + "&show=countryInfo"

		body, err := geo.httpClient.Get(url)
		if err == nil {
			var geodata GeoLocationData
			err = json.Unmarshal(body, &geodata)

			if err == nil {
				geoMeta = &GeoLocationMeta{
					GeoLocationData: &geodata,
					Error:           nil,
				}
			}
		}

		if err != nil {
			fmt.Println("Error while calling the reverse geolocation api")
			geoMeta = &GeoLocationMeta{
				GeoLocationData: nil,
				Error:           err,
			}
		}
		// update cache
		geo.cache[cacheKey] = geoMeta
		c <- geoMeta
	}
}

func NewGeoLocationResolverClient(httpClient HttpClientInterface, apiKey string) GeoLocationResolverClientInterface {
	geo := GeolocationResolver{}
	geo.apiKey = apiKey
	geo.httpClient = httpClient
	geo.cache = make(map[string]*GeoLocationMeta)

	return &geo
}
