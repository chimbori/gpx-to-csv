package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

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

	err := convertGpxFile(flag.Args())
	if err != nil {
		slog.Error("conversion failed", tint.Err(err))
		os.Exit(1)
	}
}
