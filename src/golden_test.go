package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestGoldenFiles(t *testing.T) {
	// Test the full conversion with the edam.gpx golden file
	gpxPath := filepath.Join("..", "testdata", "edam.gpx")
	goldenPath := filepath.Join("..", "testdata", "edam.csv.golden")

	// Convert to CSV
	var output bytes.Buffer
	err := convertGpxFile(&output, []string{gpxPath})
	if err != nil {
		t.Fatalf("convertGpxFile failed: %v", err)
	}

	// Read golden file
	goldenBytes, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	// Compare output with golden
	if output.String() != string(goldenBytes) {
		t.Errorf("output does not match golden file\nGot:\n%s\n\nExpected:\n%s", output.String(), string(goldenBytes))
	}
}
