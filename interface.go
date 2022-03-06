package main

type HttpClientInterface interface {
	Get(url string) ([]byte, error)
}

type CsvParserInterface interface {
	ParsePhotosCsv(fileName string) ([]*Photo, error)
}

type GeoLocationResolverClientInterface interface {
	GetGeoLocationMeta(latitude, longitude string, c chan *GeoLocationMeta)
}

type WeatherClientInterface interface {
	GetWeatherCondition(latitude, longitude, datetime string, c chan *WeatherMeta)
}

type HolidaysClientInterface interface {
	GetHolidays(year, month, day int, country, state string) (*HolidaysData, error)
}

type PhotoAlbumMetadataClientInterface interface {
	GetPhotoAlbumWithMeta([]*Photo) (*PhotoAlbum, error)
}

type ContentExaminerInterface interface {
	GetTitles(*PhotoAlbum) []string
}
