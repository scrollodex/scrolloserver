package buildit

import (
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/airtableclient"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/catutil"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/dump"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/entutil"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/locutil"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/mainlisting"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/store"
)

func FullRun(atc *airtableclient.AirClient, debugFlag bool) error {
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
	if debugFlag {
		dump.It("output-entries.yaml", mainListing)
	}

	// write individual files
	// This next line could be a goroutine.
	mainListing.WriteIndividuals()

	// write the big file
	err := mainListing.WriteBigYamlFile("data/entries.yaml")
	if err != nil {
		return err
	}

	// delete filenames no longer in use

	// Wait for any goroutines to complete.

	return nil
}
