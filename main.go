package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Printf("Usage: %s <file> <file> ...\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	parseGpx(flag.Args())
}

func parseGpx(files []string) {
	w := csv.NewWriter(os.Stdout)

	csvRow := []string{"SourceFile", "GPSDateTime", "GPSLatitude", "GPSLatitudeRef", "GPSLongitude", "GPSLongitudeRef"}
	if err := w.Write(csvRow); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, file := range files {
		log.Println(file)
		gpxBytes, err := os.ReadFile(file)
		gpx, err := ParseBytes(gpxBytes)
		if err != nil {
			log.Println(err)
		}

		for _, track := range gpx.Tracks {
			for _, segment := range track.Segments {
				for _, point := range segment.Points {
					csvRow = []string{
						"./" + point.Timestamp + ".jpg", // SourceFile,
						point.Timestamp,                 // GPSDateTime
						floatToString(point.Latitude),   // GPSLatitude
						latitudeRef(point.Latitude),     // GPSLatitudeRef
						floatToString(point.Longitude),  // GPSLongitude
						longitudeRef(point.Longitude),   // GPSLongitudeRef
					}
					if err := w.Write(csvRow); err != nil {
						log.Fatalln("error writing record to csv:", err)
					}
				}
			}
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 7, 64)
}

func latitudeRef(lat float64) string {
	if lat > 0 {
		return "North"
	} else {
		return "South"
	}
}

func longitudeRef(lon float64) string {
	if lon > 0 {
		return "East"
	} else {
		return "West"
	}
}
