package main

import (
	"encoding/json"
	"strconv"
)

type HolidaysClient struct {
	apiKey     string
	httpClient HttpClientInterface
	cache      map[string]*HolidaysData
}

func (holidays *HolidaysClient) GetHolidays(year, month, day int, country, state string) (*HolidaysData, error) {
	cacheKey := strconv.Itoa(year) + strconv.Itoa(month) + strconv.Itoa(day) + country + state
	res, ok := holidays.cache[cacheKey]
	if ok {
		return res, nil
	} else {
		var holidaysMeta *HolidaysData
		url := "https://calendarific.com/api/v2/holidays?&api_key=" + holidays.apiKey + "&country=" + country + "&location=" + state + "&year=" + strconv.Itoa(year) + "&month=" + strconv.Itoa(month) + "&day=" + strconv.Itoa(day)

		body, err := holidays.httpClient.Get(url)
		if err == nil {
			var holidayData HolidaysData
			err = json.Unmarshal(body, &holidayData)

			if err == nil {
				// update the cache
				holidays.cache[cacheKey] = &holidayData
				holidaysMeta = &holidayData
			}
		}

		return holidaysMeta, err
	}
}

func NewHolidaysClient(httpClient HttpClientInterface, apiKey string) HolidaysClientInterface {
	holidaysClient := HolidaysClient{}
	holidaysClient.apiKey = apiKey
	holidaysClient.httpClient = httpClient
	holidaysClient.cache = make(map[string]*HolidaysData)

	return &holidaysClient
}
