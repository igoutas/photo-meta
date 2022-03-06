# photo-meta

## Description
This basic command line application. allow the user to specify a photo album  (csv file of photos including timestamp and geo location) and returns a list
of suggested titles for the album. The suggested titles use metadata of each photo in the album such as Date, Weather and Holidays.

- Date: Common date entities are used, for example if all photos have been taken in a month then use the name of the month, if it's a weekend use "Weekend" ,etc.
- Weather: Retrieves historical weather data for each photos timestamp (hourly average) and the geolocation. If most photos (there is % threashold) have specific weather condition in common, then use it. E.g. "Hot", "Sunny"
- Holidays: Retrieves local and national holidays for the timestamp of each photo.

For most decisions custom logic has been applied on when an album should be characterised by a specific semantic, for example if 80% of days have rain, use Rainy. These requirements should come from the Product team.

## How to use it
Simply checkout the repository build it and call the exacutable. Then input in console the name of the file that you want to examine (the file needs to be copied under the `data` folder):
```
go build
./main
```

## What is missing
Unfortunately in this first iteration there are a lot of things that can be improved:
- Add and write unit tests (everything written as Interface already to allow easy mocking logic)
- Improve the cmd experience, like handling Interrupt events, etc.
- Log to a file rather than to the console
- Packaging. Preferably in a docker image.
- Coding style enforcement

### Future Improvements
- Introduce caching for the 3rd party API calls (memory or database) [partially done using a simple in memory map - needs to be made thread safe]. Also, e.g. for reverse geolocation do not call API when coordinates seems close enough.
- Create internal packages
- Support more timestamp formats (currently only yyyy-mm-dd hh:mm:ss is supported)
- Prepare to accept large files (parse and process csv in chunks)
- Improve the handling of cases where some photos are missing some of the Date, Holiday, Weather metadata
- Enrich the title suggestion logic by adding keyword such as "North"/"South"/"Seaside"/"Mountains" for a specific country. More personalised titles by syncing calendars and social media meta.