package locutil

import (
	"sort"
	"strings"

	"github.com/mehanizm/airtable"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/safeget"
)

// Location represents an individual geographical location.
type Location struct {
	ID          int    `yaml:"id" json:"id"`
	Display     string `yaml:"display" json:"display"`
	CountryCode string `yaml:"country_code" json:"country_code"`
	Region      string `yaml:"region" json:"region"`
	Comment     string `json:"comment,omitempty"`
}

// Locations is a list of Location.
type Locations []Location

// New allocates a new Location.
func New() Locations {
	return Locations{}
}

// ImportFromAirtable imports raw Airtable table dump.
func (store *Locations) ImportFromAirtable(raw []*airtable.Record) {
	for _, inrec := range raw {
		outrec := convert(inrec)
		*store = append(*store, *outrec)
	}
}

func convert(raw *airtable.Record) *Location {
	f := raw.Fields
	rec := &Location{
		ID:          safeget.Int(f, "x-LocationID"),
		Display:     safeget.String(f, "Location"),
		CountryCode: safeget.String(f, "x-CountryCode"),
		Region:      safeget.String(f, "x-Region"),
		Comment:     safeget.String(f, "x-Comment"),
	}
	return rec
}

// Clean or tidy the elements of the list.
func (store *Locations) Clean() {
	// Make sure the "Display" matches CC+Region+Comment.
	// Warn about any duplicates.
}

// Sort the list.
func (store *Locations) Sort() {
	// Sort with "FOO-All" at the top of all "FOO"s.

	sort.Slice(*store, func(i, j int) bool {
		return Less((*store)[i], (*store)[j])
	})
}

// Less returns a < b.
//
//	Note that "All" sorts to the top of the regions: ZA-All < ZA-Aaa < ZA-Foo
func Less(a, b Location) bool {
	if a.CountryCode != b.CountryCode {
		return a.CountryCode < b.CountryCode
	}
	if a.Region == "All" {
		return true
	}
	if b.Region == "All" {
		return false
	}
	return a.Region < b.Region
}

func (store *Locations) displays() map[string]bool {
	r := map[string]bool{}
	for _, loc := range *store {
		r[loc.Display] = true
	}
	return r
}

// ImportFromStrings takes a list of "display" locations and adds them to *store
// if they don't already exist.
func (store *Locations) ImportFromStrings(adds []string) {
	//fmt.Printf("DEBUG: Import %d locs\n", len(adds))
	locs := store.displays()
	//fmt.Printf("DEBUG: locs=%v\n", locs)
	for _, add := range adds {
		if _, ok := locs[add]; !ok {
			locs[add] = true
			store.AddLocationFromString(add)
		}
	}
	store.Sort()
}

// AddLocationFromString adds a location to store. Does not re-sort list.
func (store *Locations) AddLocationFromString(add string) {
	loc := Location{}
	loc.CountryCode, loc.Region, loc.Comment = SplitDisplay(add)
	loc.ID = store.highestID() + 1
	loc.Display = add

	//fmt.Printf("DEBUG: Adding LOCID=%d DISPLAY=%q LOC=%+v\n",
	//	loc.ID,
	//	loc.Display,
	//	loc,
	//)

	*store = append(*store, loc)
}

// highestID returns the highest ID used in *store.
func (store *Locations) highestID() int {
	max := 0
	for _, loc := range *store {
		if loc.ID > max {
			max = loc.ID
		}
	}
	return max
}

// SplitDisplay splits a Display string into the components.
func SplitDisplay(d string) (country, region, comment string) {
	d = strings.TrimSpace(d)

	spart := strings.SplitN(d, " ", 2)
	code := spart[0]
	comment = ""
	if len(spart) > 1 {
		comment = spart[1]
	}

	cpart := strings.SplitN(code, "-", 2)
	country = cpart[0]
	region = ""
	if len(cpart) > 1 {
		region = cpart[1]
	}

	if country == "ZZ" {
		country = "International"
	}

	return country, region, comment
}
