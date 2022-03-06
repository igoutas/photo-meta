package main

import (
	"time"
)

type Photo struct {
	Timestamp time.Time
	Latitude  string
	Longitude string
}

type PhotoAlbum struct {
	Photos []*PhotoMetadata
}

type PhotoMetadata struct {
	Photo        *Photo
	LocationData *GeoLocationData
	WeatherData  *WeatherData
	HolidaysData *HolidaysData
	IsWeekend    bool
}

type GeoLocationData struct {
	Items []struct {
		Title   string `json:"title"`
		Type    string `json:"place"`
		Address struct {
			FullAddress string `json:"label"`
			Country     string `json:"countryName"`
			State       string `json:"state"`
			City        string `json:"city"`
			Street      string `json:"street"`
			District    string `json:"district"`
			County      string `json:"county"`
		} `json:"address"`
		Labels []struct {
			Name string `json:"name"`
		} `json:"categories"`
		CountryInfo struct {
			Alpha2CountryCode string `json:"alpha2"`
			Alpha3CountryCode string `json:"alpha3"`
		} `json:"countryInfo"`
	} `json:"items"`
}

type GeoLocationMeta struct {
	GeoLocationData *GeoLocationData
	Error           error
}

type WeatherData struct {
	Days []struct {
		Conditions  string `json:"conditions"`
		Description string `json:"description"`
	}
	HourConditions struct {
		Temperature float32 `json:"temp"`
		CloudCover  float32 `json:"cloudcover"`
		Snow        float32 `json:"snow"`
		Conditions  string  `json:"conditions"`
		Visibility  float32 `json:"visibility"`
	} `json:"currentConditions"`
}

type WeatherMeta struct {
	WeatherData *WeatherData
	Error       error
}

type HolidaysData struct {
	Response struct {
		Holidays []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Country     struct {
				Name string `json:"name"`
			} `json:"country"`
			HolidayType []string `json:"type"`
		} `json:"holidays"`
	} `json:"response"`
}
