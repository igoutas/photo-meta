package main

import (
	"strconv"
	"strings"
)

type ContentExaminer struct {
}

type DateSemantics struct {
	Day     string
	Weekend string
	Week    string
	Month   string
	Season  string
	Year    string
}

func (contentExaminer *ContentExaminer) GetTitles(album *PhotoAlbum) []string {
	var titles, holidaysSemantics []string
	locationSemantics := contentExaminer.getLocationSemantics(album)
	weatherSemantics := contentExaminer.getWeatherSemantics(album)
	dateSemantics, shouldExamineHolidays := contentExaminer.getDateSemantics(album)

	if shouldExamineHolidays {
		holidaysSemantics = contentExaminer.getHolidaysSemantics(album)
	}

	// Construct titles based on templates
	titles = append(titles, contentExaminer.getDateLocationTitles(locationSemantics, dateSemantics)...)
	titles = append(titles, contentExaminer.getWeatherDateLocationTitles(locationSemantics, weatherSemantics, dateSemantics)...)
	titles = append(titles, contentExaminer.getWeatherLocationTitles(locationSemantics, weatherSemantics)...)
	titles = append(titles, contentExaminer.getHolidayTitles(holidaysSemantics, dateSemantics)...)
	titles = append(titles, contentExaminer.getLocationTitles(locationSemantics)...)

	return titles
}

func (contentExaminer *ContentExaminer) getLocationTitles(locationSemantics []string) []string {
	var titles []string
	// Location specific titles
	if len(locationSemantics) > 0 {
		for _, location := range locationSemantics {
			titles = append(titles, location+" Memories")
		}
	}

	return titles
}

func (contentExaminer *ContentExaminer) getHolidayTitles(holidaysSemantics []string, dateSemantics *DateSemantics) []string {
	var titles []string
	// If there is a Holiday detected associate it with the day or the weekend if exists
	if len(holidaysSemantics) > 0 {
		for _, holiday := range holidaysSemantics {
			if len(dateSemantics.Day) > 0 {
				titles = append(titles, "Celebrating "+holiday)
			}
			if len(dateSemantics.Weekend) > 0 {
				titles = append(titles, holiday+" Weekend")
			}
		}
	}

	return titles
}

func (contentExaminer *ContentExaminer) getWeatherLocationTitles(locationSemantics, weatherSemantics []string) []string {
	var titles []string
	// {weather} {place}. e.g. Rainy New York
	if len(weatherSemantics) > 0 && len(locationSemantics) > 0 {
		for _, weather := range weatherSemantics {
			for _, location := range locationSemantics {
				titles = append(titles, weather+" "+location)
			}
		}
	}

	return titles
}

func (contentExaminer *ContentExaminer) getWeatherDateLocationTitles(locationSemantics, weatherSemantics []string, dateSemantics *DateSemantics) []string {
	var titles []string
	// A {weather} {Date} in {location}
	if len(weatherSemantics) > 0 && len(locationSemantics) > 0 {
		for _, weather := range weatherSemantics {
			for _, location := range locationSemantics {
				if len(dateSemantics.Day) > 0 {
					titles = append(titles, "A "+weather+" "+dateSemantics.Day+" in "+location)
				}
				if len(dateSemantics.Weekend) > 0 {
					titles = append(titles, "A "+weather+" "+dateSemantics.Weekend+" in "+location)
				}
				if len(dateSemantics.Week) > 0 {
					titles = append(titles, "A "+weather+" "+dateSemantics.Week+" in "+location)
				}
				if len(dateSemantics.Month) > 0 {
					titles = append(titles, "A "+weather+" "+dateSemantics.Month+" in "+location)
				}
				if len(dateSemantics.Season) > 0 {
					titles = append(titles, "A "+weather+" "+dateSemantics.Season+" in "+location)
				}
				if len(dateSemantics.Year) > 0 {
					titles = append(titles, "A "+weather+" "+dateSemantics.Year+" in "+location)
				}
			}
		}
	}

	return titles
}

func (contentExaminer *ContentExaminer) getDateLocationTitles(locationSemantics []string, dateSemantics *DateSemantics) []string {
	var titles []string
	// e.g. {date} in {place}
	if len(locationSemantics) > 0 {
		for _, location := range locationSemantics {
			if len(dateSemantics.Day) > 0 {
				titles = append(titles, "A "+dateSemantics.Day+" in "+location)
			}

			if len(dateSemantics.Week) > 0 {
				titles = append(titles, "A "+dateSemantics.Week+" in "+location)
			}

			if len(dateSemantics.Month) > 0 {
				titles = append(titles, "A "+dateSemantics.Month+" in "+location)
			}

			if len(dateSemantics.Year) > 0 {
				titles = append(titles, "A "+dateSemantics.Year+" in "+location)
			}

			if len(dateSemantics.Season) > 0 {
				titles = append(titles, location+" in "+dateSemantics.Season)
			}
		}
	}

	//If no specific location use something like Here & There 2019
	if len(locationSemantics) == 0 && len(dateSemantics.Year) > 0 {
		titles = append(titles, "Here & There "+dateSemantics.Year)
	}

	return titles
}

func (contentExaminer *ContentExaminer) getHolidaysSemantics(album *PhotoAlbum) []string {
	var holidaysSemantics map[string]struct{} // in order to avoid duplicates
	var holidays []string
	var empty struct{}

	for _, photo := range album.Photos {
		for _, holiday := range photo.HolidaysData.Response.Holidays {
			if contains(holiday.HolidayType, "National holiday") || contains(holiday.HolidayType, "Local holiday") {
				holidaysSemantics[holiday.Name] = empty
			}
		}
	}

	for k := range holidaysSemantics {
		holidays = append(holidays, k)
	}

	return holidays
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (contentExaminer *ContentExaminer) getLocationSemantics(album *PhotoAlbum) []string {
	// try to find the location types that all photos have in common, e.g. city, district, etc.
	var locationSemantics []string
	hasSameTitle, hasSameDistrict, hasSameCity, hasSameCountry, hasSameCounty := true, true, true, true, true
	city := album.Photos[0].LocationData.Items[0].Address.City
	district := album.Photos[0].LocationData.Items[0].Address.District
	country := album.Photos[0].LocationData.Items[0].Address.Country
	county := album.Photos[0].LocationData.Items[0].Address.County
	locationTitle := album.Photos[0].LocationData.Items[0].Title

	for _, photo := range album.Photos {
		if len(photo.LocationData.Items) == 0 {
			continue // we could possibly return an error here, up for discussion
		}

		if hasSameTitle {
			hasSameTitle = photo.LocationData.Items[0].Title == locationTitle
		}

		if hasSameDistrict {
			hasSameDistrict = photo.LocationData.Items[0].Address.District == district
		}

		if hasSameCity {
			hasSameCity = photo.LocationData.Items[0].Address.City == city
		}

		if hasSameCountry {
			hasSameCountry = photo.LocationData.Items[0].Address.Country == country
		}

		if hasSameCounty {
			hasSameCounty = photo.LocationData.Items[0].Address.County == county
		}
	}

	if hasSameTitle {
		locationSemantics = append(locationSemantics, locationTitle)
	}

	if hasSameDistrict {
		locationSemantics = append(locationSemantics, district)
	}

	if hasSameCity {
		locationSemantics = append(locationSemantics, city)
	}

	if hasSameCountry {
		locationSemantics = append(locationSemantics, country)
	}

	if hasSameCounty {
		locationSemantics = append(locationSemantics, county)
	}

	return locationSemantics
}

func (contentExaminer *ContentExaminer) getDateSemantics(album *PhotoAlbum) (*DateSemantics, bool) {
	// try to find a Date related description of the alnum, e.g. Week, Weekend, Month, Year, etc,
	dateSemantics := DateSemantics{}
	isWeekend, hasSameSeason := true, true
	minDate, maxDate := album.Photos[0].Photo.Timestamp, album.Photos[0].Photo.Timestamp

	for _, photo := range album.Photos {
		if photo.Photo.Timestamp.Unix() < minDate.Unix() {
			minDate = photo.Photo.Timestamp
		}

		if photo.Photo.Timestamp.Unix() > maxDate.Unix() {
			maxDate = photo.Photo.Timestamp
		}

		if hasSameSeason {
			hasSameSeason = (int(minDate.Month())-1)/3 == (int(photo.Photo.Timestamp.Month())-1)/3
		}

		// Check if all photos were from the same weekend
		if isWeekend {
			isWeekend = photo.IsWeekend
		}
	}

	numOfDays := maxDate.Sub(minDate).Hours() / 24
	if minDate.Day() == maxDate.Day() {
		dateSemantics.Day = "Day"
	}

	// Consider it a week if it's between 5-7 days
	if numOfDays < 8.0 && numOfDays > 4.0 {
		dateSemantics.Week = "Week"
	}

	// Consider it a Month if it's between 25 - 34 days
	if numOfDays < 35.0 && numOfDays > 24.0 {
		if minDate.Month() != maxDate.Month() {
			dateSemantics.Month = "Month"
		} else {
			// All photos belong to the same Month
			dateSemantics.Month = minDate.Month().String()
		}
	}

	if hasSameSeason {
		var season string
		lat, _ := strconv.ParseFloat(album.Photos[0].Photo.Latitude, 32)
		isNorthHemisphere := lat > 0

		if (int(minDate.Month()))/3 == 0 || int(minDate.Month()) == 12 {
			if isNorthHemisphere {
				season = "Winter"
			} else {
				season = "Summer"
			}
		} else if (int(minDate.Month()))/3 == 1 {
			if isNorthHemisphere {
				season = "Spring"
			} else {
				season = "Autumn"
			}
		} else if (int(minDate.Month()))/3 == 2 || int(minDate.Month()) == 12 {
			if isNorthHemisphere {
				season = "Summer"
			} else {
				season = "Winter"
			}
		} else {
			// 9,10, 11 months
			if isNorthHemisphere {
				season = "Autumn"
			} else {
				season = "Spring"
			}
		}
		dateSemantics.Season = season
	}

	// If photos are scattered in more than 10 months append the Year label
	if numOfDays > 310.0 {
		if minDate.Year() == maxDate.Year() {
			dateSemantics.Year = strconv.Itoa(minDate.Year())
		} else {
			dateSemantics.Year = "Year"
		}
	}

	if isWeekend && numOfDays < 3.0 {
		dateSemantics.Weekend = "Weekend"
	}

	// Associate with specific hliday if it's a one day trip or a weekend
	return &dateSemantics, isWeekend || numOfDays == 1
}

func (contentExaminer *ContentExaminer) getWeatherSemantics(album *PhotoAlbum) []string {
	// gather weather semantics. If description present in more than 80% of the photos use it.
	var hot, cloudy, rainy, sunny, snowy, foggy int
	var weatherSemantics []string

	for _, photo := range album.Photos {
		// if temperature more than 25 Celcius then consider it hot
		if photo.WeatherData.HourConditions.Temperature > 25 {
			hot += 1
		}

		// if Snow exists conmsider it snowy
		if photo.WeatherData.HourConditions.Snow > 0.1 {
			snowy += 1
		}

		// if can't see further that 0.2km then consider it foggy
		if photo.WeatherData.HourConditions.Visibility < 0.2 {
			foggy += 1
		}

		// If clou coverage > 87% consider it cloudy
		if photo.WeatherData.HourConditions.CloudCover > 87 {
			cloudy += 1
		}

		// if rain description was present for the current hour of the photo
		if strings.Contains(photo.WeatherData.HourConditions.Conditions, "Rain") {
			rainy += 1
		}

		// if Clear description was present for the current hour of the photo
		if strings.Contains(photo.WeatherData.HourConditions.Conditions, "Clear") {
			sunny += 1
		}

	}

	if float32(hot*100/len(album.Photos)) > float32(80) {
		weatherSemantics = append(weatherSemantics, "Hot")
	}

	if float32(snowy*100/len(album.Photos)) > float32(80) {
		weatherSemantics = append(weatherSemantics, "Snowy")
	}

	if float32(foggy*100/len(album.Photos)) > float32(80) {
		weatherSemantics = append(weatherSemantics, "Foggy")
	}

	if float32(cloudy*100/len(album.Photos)) > float32(80) {
		weatherSemantics = append(weatherSemantics, "Cloudy")
	}

	if float32(rainy*100/len(album.Photos)) > float32(80) {
		weatherSemantics = append(weatherSemantics, "Rainy")
	}

	if float32(sunny*100/len(album.Photos)) > float32(80) {
		weatherSemantics = append(weatherSemantics, "Sunny")
	}

	return weatherSemantics
}

func NewContentExaminer() ContentExaminerInterface {
	return &ContentExaminer{}
}
