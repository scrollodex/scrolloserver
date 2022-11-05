package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/airtableclient"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/buildit"
)

var debugFlag bool

func main() {
	flag.BoolVar(&debugFlag, "debug", false, "Output debug info")
	flag.Parse()

	fmt.Println("air2hugo running!")

	// TODO(tlim): Verify env variables are set.

	atc := airtableclient.New(
		os.Getenv("AIRTABLE_APIKEY"),
		os.Getenv("AIRTABLE_BASE_ID"),
	)

	err := buildit.FullRun(atc, debugFlag)
	if err != nil {
		log.Fatal(err)
	}
}
