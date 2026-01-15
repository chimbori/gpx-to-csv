package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/lmittmann/tint"
)

func convertGpxFile(output io.Writer, files []string) error {
	w := csv.NewWriter(output)
	defer w.Flush()

	csvRow := []string{"SourceFile", "GPSDateTime", "GPSLatitude", "GPSLatitudeRef", "GPSLongitude", "GPSLongitudeRef"}
	if err := w.Write(csvRow); err != nil {
		slog.Error("error writing record to csv", tint.Err(err))
		os.Exit(1)
	}

	for _, file := range files {
		slog.Info(file)

		gpxBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		gpx, err := parseGpx(gpxBytes)
		if err != nil {
			return fmt.Errorf("error parsing file: %w", err)
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
						return fmt.Errorf("error writing CSV: %w", err)
					}
				}
			}
		}
	}

	return nil
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
