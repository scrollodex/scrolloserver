package entutil

import (
	"fmt"
	"sort"

	"github.com/mehanizm/airtable"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/locutil"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/safeget"
)

// Entry represents an individual resource provider.
type Entry struct {
	Title string `yaml:"title" json:"title"`
	ID    int    `yaml:"id" json:"id"`

	Salutation  string `yaml:"salutation" json:"salutation"`
	Firstname   string `yaml:"first_name" json:"first_name"`
	Lastname    string `yaml:"last_name" json:"last_name"`
	Credentials string `yaml:"credentials" json:"credentials"`
	JobTitle    string `yaml:"job_title" json:"job_title"`

	Company string `yaml:"company" json:"company"`

	ShortDesc string `yaml:"short_desc" json:"short_desc"` // MarkDown (1 line)
	Phone     string `yaml:"phone" json:"phone"`
	//Fax       string `yaml:"fax" json:"fax"`
	Address string `yaml:"address" json:"address"`

	Email  string `yaml:"email" json:"email"`
	Email2 string `yaml:"email2" json:"email2"`

	Website  string `yaml:"website" json:"website"`
	Website2 string `yaml:"website2" json:"website2"`

	Fees        string `yaml:"fees" json:"fees"`               // MarkDown
	Description string `yaml:"description" json:"description"` // MarkDown

	Category     string   `yaml:"categories" json:"categories"`
	Location     []string `json:"location"`
	Country      string   `yaml:"countries" json:"countries,omitempty"`
	Region       string   `yaml:"regions" json:"regions"`
	Status       int      `yaml:"-" json:"status"` // 0=Inactive, 1=Active, 2=Proposed
	LastEditDate string   `yaml:"last_update" json:"last_update"`
}

// Entries is a list of Entry.
type Entries []Entry

// New allocates a new Entry.
func New() Entries {
	return Entries{}
}

// ImportFromAirtable imports raw Airtable table dump.
func (store *Entries) ImportFromAirtable(raw []*airtable.Record) {
	for _, inrec := range raw {
		outrec := convert(inrec)
		if outrec != nil {
			*store = append(*store, *outrec)
		}
	}
}

func convert(raw *airtable.Record) *Entry {
	f := raw.Fields

	lastmod := safeget.String(f, "Last Modified")
	if lastmod == "" {
		lastmod = safeget.String(f, "x-lastUpdate")
	}
	if len(lastmod) > 10 {
		lastmod = lastmod[:10]
	}

	rec := &Entry{
		ID:       getID(f),
		Category: safeget.String(f, "Category"),
		Location: safeget.Strings(f, "Location"),
		Status:   getStatus(f),

		Company:     safeget.String(f, "Company"),
		Salutation:  safeget.String(f, "Sal"),
		Firstname:   safeget.String(f, "First"),
		Lastname:    safeget.String(f, "Last"),
		Credentials: safeget.String(f, "Suffix"),
		JobTitle:    safeget.String(f, "Job_Title"),

		ShortDesc:   safeget.String(f, "Short Description"),
		Description: safeget.String(f, "Description"),
		Fees:        safeget.String(f, "Fees"),

		Address: safeget.String(f, "Address"),
		Email:   safeget.String(f, "Email"),
		Email2:  safeget.String(f, "Email2"),
		Phone:   safeget.String(f, "Phone"),
		//Fax:      safeget.String(f, "Fax"),
		Website:  safeget.String(f, "Website"),
		Website2: safeget.String(f, "Website2"),

		LastEditDate: lastmod,
	}

	// FIXME: only deals with the first location
	if len(rec.Location) == 0 {
		rec.Location = []string{"Unknown"}
	}
	c, r, _ := locutil.SplitDisplay(rec.Location[0])
	rec.Country = c
	rec.Region = r

	if rec.Category == "" {
		rec.Category = "unknown"
		fmt.Printf("UNKNOWN CATEGORY: %+v\n", f)
	}

	return rec
}

func getID(f map[string]interface{}) int {
	// The ID is the "EntryID" column or, if blank,
	// use the auto-generated ID column.

	r := safeget.Int(f, "EntryID")
	if r != 0 {
		return r
	}
	// No Legacy ID?  Use the autonumber id.
	//fmt.Printf("Autonumber %v is new-style (no legacy EntryID)\n", entryID)
	return safeget.Int(f, "ID")
}

func getStatus(f map[string]interface{}) int {
	switch v := f["Status"].(type) {
	case string:
		switch v {
		case "HIDDEN":
			return 0
		case "SHOW":
			return 1
		case "PROPOSED":
			return 2
		}
	default:
		fmt.Printf("DEBUG: UNKNOWN STATUS %v = %v\n", v, f["Status"])
	}
	return 99
}

// Clean or tidy the elements of the list.
func (store *Entries) Clean() {
}

// Sort the list.
func (store *Entries) Sort() {
	// FIXME: This might not be needed. I think Hugo sorts
	// them for us.  At least for now, sort them to match the order of
	// the legacy system so we have an easier to comparing the outputs.

	sort.Slice(*store, func(i, j int) bool {
		return Less((*store)[i], (*store)[j])
	})

}

// Less sorts on ID.
func Less(a, b Entry) bool {
	return a.ID < b.ID
}

// Categories returns a list of the category names used
// in store.  The list is de-duped and sorted.
func (store *Entries) Categories() []string {
	var result []string

	seen := map[string]bool{}
	for _, item := range *store {
		n := item.Category
		if !seen[n] {
			result = append(result, n)
			seen[n] = true
		}
	}

	sort.Strings(result)

	return result
}

// Locations returns a list of the country "display strings" used
// in store.  The list is de-duped and sorted.
func (store *Entries) Locations() []string {
	var result []string

	seen := map[string]bool{}
	for _, item := range *store {
		for _, n := range item.Location {
			if !seen[n] {
				result = append(result, n)
				seen[n] = true
			}
		}
	}

	sort.Strings(result)

	return result
}

// // FlattenEntriesOnePerLocation takes a list of entries and returns a new list
// // with each entry duplicated for each of its locations.
// func FlattenEntriesOnePerLocation(ents Entries) Entries {
// 	var result Entries
// 	for _, ent := range ents {
// 		if len(ent.Location) == 0 {
// 			panic("entry has no location") // FIXME: Handle this case more gracefully.  entutil.go should have assured that all entries have at least one location (even if it is "Unknown")
// 		}
// 		for _, loc := range ent.Location {
// 			// Create a copy of the entry with the current location
// 			newEnt := ent
// 			newEnt.Location = []string{loc}
// 			result = append(result, newEnt)
// 		}
// 	}
// 	return result
// }
