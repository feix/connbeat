package main

import (
	"os"

	_ "github.com/raboof/beats-output-http"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/raboof/connbeat/beater"
)

var Name = "connbeat"

func main() {
	if err := beat.Run(Name, appVersion, beater.New); err != nil {
		os.Exit(1)
	}
}
