package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	hereApiKey           string = "nQJvCn3_fF_DZ_Qfaa-JjHXoGvVwfGXTIt2mDWIv3No"
	visualcrossingApiKey string = "Z53BNDBXAZ4REX94ESGJBDWE7"
	calendarificApiKey   string = "16aed471fbcc81dc9c2ced9a7f2c7188628cbe3f"
)

func main() {
	datapath := getDataPath()
	parser := NewCsvParser(datapath)
	dataFiles := getDataFiles(datapath)
	httpClient := NewHttpClient()
	geoLocationResolverClient := NewGeoLocationResolverClient(httpClient, hereApiKey)
	weatherClient := NewWeatherClient(httpClient, visualcrossingApiKey)
	holidaysClient := NewHolidaysClient(httpClient, calendarificApiKey)
	photoAlbumMetaClient := NewPhotoAlbumMetadataClient(geoLocationResolverClient, weatherClient, holidaysClient)
	contentExaminer := NewContentExaminer()

	// List the available data files to the user
	fmt.Printf("Please specify (console input) any of the available data files to process:\n %v \n", dataFiles)

	reader := bufio.NewReader(os.Stdin)
	dataFile, err := reader.ReadString('\n')

	if err != nil {
		panic(err)
	}

	photos, err := parser.ParsePhotosCsv(strings.TrimSuffix(dataFile, "\n"))

	if err != nil {
		panic(err)
	}

	start := time.Now()
	album, err := photoAlbumMetaClient.GetPhotoAlbumWithMeta(photos)
	fmt.Printf("Photo meta took %v\n", time.Since(start))
	if err != nil {
		panic(err)
	}

	titles := contentExaminer.GetTitles(album)
	output := "'" + strings.Join(titles, `','`) + `'`

	fmt.Println("----------------------------------")
	fmt.Println("Suggested titles for: " + dataFile)
	fmt.Println(output)
}

func getDataPath() string {
	ex, err := os.Executable()
	if err != nil {
		fmt.Println("Error while locating the current execution path")
		panic(err)
	}
	exPath := filepath.Dir(ex)
	dataPath := filepath.Join(exPath, "/data")

	_, err = os.Stat(dataPath)
	if err == nil {
		return dataPath
	}

	if err != nil {
		fmt.Println("Error while locating data file path")
		panic(err)
	}

	return dataPath
}

func getDataFiles(directory string) []string {
	var dataFiles []string
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		dataFiles = append(dataFiles, f.Name())
	}

	return dataFiles
}
