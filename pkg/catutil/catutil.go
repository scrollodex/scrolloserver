package catutil

import (
	"sort"

	"github.com/mehanizm/airtable"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/safeget"
)

// Category represents an type of service provided.
type Category struct {
	ID          int    `yaml:"id" json:"id"`
	Name        string `yaml:"category_name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Priority    int    `yaml:"priority" json:"priority"`
	Icon        string `yaml:"icon" json:"icon"`
}

// Categories is a list of Category.
type Categories []Category

// New allocates a new Category.
func New() Categories {
	return Categories{}
}

// ImportFromAirtable imports raw Airtable table dump.
func (store *Categories) ImportFromAirtable(raw []*airtable.Record) {
	for _, inrec := range raw {
		//fmt.Printf("DEBUG: category = IN  %+v\n", inrec)
		outrec := convert(inrec)
		//fmt.Printf("DEBUG: category = OUT %+v\n\n", outrec)
		*store = append(*store, *outrec)
	}
}

func convert(raw *airtable.Record) *Category {
	f := raw.Fields
	xcatid := safeget.Int(f, "x-CategoryID")

	rec := &Category{
		ID:          xcatid,
		Name:        safeget.String(f, "Name"),
		Description: safeget.String(f, "Description"),
		Priority:    safeget.Int(f, "x-Priority"),
		Icon:        safeget.String(f, "IconFilename"),
	}
	return rec
}

// Clean or tidy the elements of the list.
func (store *Categories) Clean() {
	store.fixIDs()
}

// Sort the list.
func (store *Categories) Sort() {
	// Sort by priority, then name.

	sort.Slice(*store, func(i, j int) bool {
		return Less((*store)[i], (*store)[j])
	})
}

// Less returns a < b.
//
//	Sort by priority, then name.
func Less(a, b Category) bool {
	ap := a.Priority
	bp := b.Priority

	// No priority?  Leave it for last.
	if ap == 0 {
		ap = 99999
	}
	if bp == 0 {
		bp = 99999
	}

	if ap != bp {
		return ap < bp
	}
	return a.Name < b.Name
}

func (store *Categories) categories() map[string]bool {
	r := map[string]bool{}
	for _, loc := range *store {
		r[loc.Name] = true
	}
	return r
}

// ImportFromStrings takes a list of categories and adds them to *store
// if they don't already exist.
func (store *Categories) ImportFromStrings(adds []string) {
	//fmt.Printf("DEBUG: Import %d cats\n", len(adds))
	locs := store.categories()
	for _, add := range adds {
		if _, ok := locs[add]; !ok {
			locs[add] = true
			store.AddCategoryFromString(add)
		}
	}
	store.Sort()
}

// AddCategoryFromString adds a category to store. Does not re-sort list.
func (store *Categories) AddCategoryFromString(add string) {
	cat := Category{}
	cat.Name = add
	cat.ID = store.highestID() + 1
	cat.Priority = 1

	//fmt.Printf("DEBUG: Adding CATID=%d NAME=%q CAT=%+v\n",
	//	cat.ID,
	//	cat.Name,
	//	cat,
	//)

	*store = append(*store, cat)
}

// highestID returns the highest ID used in *store.
func (store *Categories) highestID() int {
	max := 0
	for _, cat := range *store {
		if cat.ID > max {
			max = cat.ID
		}
	}
	return max
}

// fixMissingIDs fixes the ID field.  The goal is to remove any
// duplicates, fill in any missing (zero) ids.  The easiest way to do
// this is to just sort the list by name and renumber the items.
// This is fine because no system uses the IDs at this point.
func (store *Categories) fixIDs() {
	store.Sort()
	for i, _ := range *store {
		(*store)[i].ID = i + 1
	}
}
