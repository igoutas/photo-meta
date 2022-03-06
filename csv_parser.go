package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type CsvParser struct {
	dataPath string
}

func (parser *CsvParser) ParsePhotosCsv(fileName string) ([]*Photo, error) {

	filePath := filepath.Join(parser.dataPath, fileName)
	csvfile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Unable to open the data file", err)
		return nil, err
	}

	var photos []*Photo
	r := csv.NewReader(csvfile)

	// Iterate through the records
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error while reading file")
			return nil, err
		}
		// Construct the metadata
		timestamp, err := time.Parse("2006-01-02 15:04:05", record[0])
		if err != nil {
			fmt.Println("Malformed date found")
			return nil, err
		}

		// Check that latitude and longitude have valid float values, otherwise subsequent geo location calls will fail
		_, err = strconv.ParseFloat(record[1], 64)
		if err != nil {
			fmt.Println("Malformed latitude found")
			return nil, err
		}

		_, err = strconv.ParseFloat(record[2], 64)
		if err != nil {
			fmt.Println("Malformed longitude found")
			return nil, err
		}

		photo := Photo{
			Timestamp: timestamp,
			Latitude:  record[1],
			Longitude: record[2],
		}
		photos = append(photos, &photo)
	}

	return photos, nil
}

func NewCsvParser(dataPath string) CsvParserInterface {
	parser := CsvParser{}

	parser.dataPath = dataPath

	return &parser
}
