package main

import (
	"os"

	"github.com/edrans/cloudflarebeat/beater"

	"github.com/elastic/beats/libbeat/beat"
)

func main() {
	err := beat.Run("cloudflarebeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
