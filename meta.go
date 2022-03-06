package main

import (
	"fmt"
	"time"
)

type PhotoAlbumMetadataClient struct {
	geoLocationResolverClient GeoLocationResolverClientInterface
	weatherClient             WeatherClientInterface
	holidaysClient            HolidaysClientInterface
}

func (photosAlbum *PhotoAlbumMetadataClient) GetPhotoAlbumWithMeta(photos []*Photo) (*PhotoAlbum, error) {
	var photosMeta []*PhotoMetadata
	geoMetaChan := make(chan *GeoLocationMeta)
	weatherMetaChan := make(chan *WeatherMeta)

	for _, photo := range photos {
		var geoMeta *GeoLocationData
		var holidaysMeta *HolidaysData
		var weatherMeta *WeatherData
		var err error

		// TODO create a pool of threads in order to avoid thread construction/destruction latency
		go photosAlbum.geoLocationResolverClient.GetGeoLocationMeta(photo.Latitude, photo.Longitude, geoMetaChan)
		go photosAlbum.weatherClient.GetWeatherCondition(photo.Latitude, photo.Longitude, photo.Timestamp.Format(time.RFC3339), weatherMetaChan)

		for i := 0; i < 2; i++ {
			select {
			case geo := <-geoMetaChan:
				if geo.Error != nil {
					return nil, geo.Error
				}
				geoMeta = geo.GeoLocationData
			case weather := <-weatherMetaChan:
				if weather.Error != nil {
					return nil, weather.Error
				}
				weatherMeta = weather.WeatherData
			}
		}

		if len(geoMeta.Items) > 0 {
			countryCodeAlpha2 := geoMeta.Items[0].CountryInfo.Alpha2CountryCode
			state := geoMeta.Items[0].Address.State
			holidaysMeta, err = photosAlbum.holidaysClient.GetHolidays(photo.Timestamp.Year(), (int)(photo.Timestamp.Month()), photo.Timestamp.Day(), countryCodeAlpha2, state)

			if err != nil {
				fmt.Println("Error while retrieving holidays data for photo")
				return nil, err
			}
		}

		photoMeta := PhotoMetadata{
			Photo:        photo,
			LocationData: geoMeta,
			WeatherData:  weatherMeta,
			HolidaysData: holidaysMeta,
			IsWeekend:    int(photo.Timestamp.Weekday()) > 5,
		}

		photosMeta = append(photosMeta, &photoMeta)
	}

	album := PhotoAlbum{Photos: photosMeta}

	return &album, nil
}

func NewPhotoAlbumMetadataClient(geoLocationResolverClient GeoLocationResolverClientInterface, weatherClient WeatherClientInterface, holidaysClient HolidaysClientInterface) PhotoAlbumMetadataClientInterface {
	album := PhotoAlbumMetadataClient{}

	album.geoLocationResolverClient = geoLocationResolverClient
	album.weatherClient = weatherClient
	album.holidaysClient = holidaysClient

	return &album
}
