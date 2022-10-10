package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/airtableclient"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/catutil"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/dump"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/entutil"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/locutil"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/mainlisting"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/store"
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

	// go GatherExistingEntryFilenames()

	// go rawLocs.PullTable()
	//    locData.FromAirtable(rawLocs)
	//    locData.Clean()

	locations := locutil.New()
	categories := catutil.New()
	entries := entutil.New()
	store.Setup(atc, &locations, "Locations")
	store.Setup(atc, &categories, "Categories")
	store.Setup(atc, &entries, "Entries")

	// wait for pulls to complete.

	// Find any missing categories and add them.
	categories.ImportFromStrings(entries.Categories())
	// Find any missing locations and add them.
	locations.ImportFromStrings(entries.Locations())

	if debugFlag {
		dump.It("final-categories.json", categories)
		dump.It("final-locations.json", locations)
		dump.It("final-entries.json", entries)
	}

	mainListing := mainlisting.Assemble(categories, locations, entries)
	//dump.It("output-entries.yaml", mainListing)

	// This next line can be a goroutine.
	// write individual files
	mainListing.WriteIndividuals()

	// write the big file
	err := mainListing.WriteBigYamlFile("data/entries.yaml")
	if err != nil {
		log.Fatal(err)
	}
	// delete filenames no longer in use

	// Wait for any goroutines to complete.
}
