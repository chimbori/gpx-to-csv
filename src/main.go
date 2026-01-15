package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/lmittmann/tint"
)

func main() {
	tintHandler := tint.NewHandler(os.Stderr, &tint.Options{TimeFormat: "2006-01-02 15:04:05.000"})
	slog.SetDefault(slog.New(tintHandler))

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
		slog.Error("error writing record to csv", tint.Err(err))
		os.Exit(1)
	}

	for _, file := range files {
		slog.Info(file)
		gpxBytes, err := os.ReadFile(file)
		gpx, err := ParseBytes(gpxBytes)
		if err != nil {
			slog.Error("parse error", tint.Err(err))
		}

		for _, track := range gpx.Tracks {
			for _, segment := range track.Segments {
				for _, point := range segment.Points {
					localTimestamp := utcToLocal(point.Timestamp)
					csvRow = []string{
						"./" + localTimestamp + ".jpg",   // SourceFile,
						localTimestamp,                   // GPSDateTime
						precision7digit(point.Latitude),  // GPSLatitude
						latitudeRef(point.Latitude),      // GPSLatitudeRef
						precision7digit(point.Longitude), // GPSLongitude
						longitudeRef(point.Longitude),    // GPSLongitudeRef
					}
					if err := w.Write(csvRow); err != nil {
						slog.Error("error writing record to csv", tint.Err(err))
						os.Exit(1)
					}
				}
			}
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		slog.Error("csv flush error", tint.Err(err))
		os.Exit(1)
	}
}

func utcToLocal(utcTimeStr string) string {
	utcTime, err := time.Parse(time.RFC3339, utcTimeStr)
	if err != nil {
		slog.Error("error parsing timestamp", tint.Err(err), "timestamp", utcTimeStr)
		return utcTimeStr // Return original if parsing fails.
	}
	return utcTime.Local().Format(time.RFC3339)
}

func precision7digit(f float64) string { return strconv.FormatFloat(f, 'f', 7, 64) }

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
