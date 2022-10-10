package airtableclient

import (
	"fmt"

	"github.com/mehanizm/airtable" // https://pkg.go.dev/github.com/mehanizm/airtable#section-readme
)

// AirClient is the handle for accessing an Airtable account.
type AirClient struct {
	client *airtable.Client
	apikey string
	baseid string
}

// New allocates a new AirClient.
func New(apikey, baseid string) *AirClient {
	r := AirClient{
		apikey: apikey,
		baseid: baseid,
	}
	r.client = airtable.NewClient(r.apikey)
	return &r
}

// GetRecordsAll returns the raw records from a table by its name.
func (atc *AirClient) GetRecordsAll(tableName string) ([]*airtable.Record, error) {

	// Look up the table's ID.
	tab := atc.client.GetTable(atc.baseid, tableName)

	// Get the records:
	raw, err := downloadTableRecords(tab)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s's records: %w", tableName, err)
	}

	return raw, nil
}

// downloadTableRecords gets all records of a table (all pages).
func downloadTableRecords(table *airtable.Table) ([]*airtable.Record, error) {
	// TODO(tlim): If we are rate limited, retry.

	var result []*airtable.Record

	var offset string
	for {
		// Get 1 page of records.
		records, err := table.GetRecords().
			WithOffset(offset).
			Do()
		if err != nil {
			return nil, err
		}
		result = append(result, records.Records...)

		// Stop when we're out of records.
		offset = records.Offset
		if offset == "" {
			break
		}
	}

	return result, nil
}
