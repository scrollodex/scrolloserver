package store

import (
	"github.com/mehanizm/airtable"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/airtableclient"
)

// Store defines an interface that stores a table of data.
type Store interface {
	ImportFromAirtable([]*airtable.Record)
	Clean()
	Sort()
}

// Setup pulls the data from Airtable and returns cleaned, sorted, converted Store.
func Setup(atc *airtableclient.AirClient, store Store, tableName string) {
	raw, err := atc.GetRecordsAll(tableName)
	if err != nil {
		panic(err)
	}

	//dump.It("raw-"+tableName+".json", raw)

	store.ImportFromAirtable(raw)
	store.Clean()
	store.Sort()
}
