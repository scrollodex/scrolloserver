package mainlisting

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/catutil"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/entutil"
	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/locutil"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"gopkg.in/yaml.v3"
)

// MainListing represents the yaml output file delivered to Hugo.
type MainListing struct {
	Categories     catutil.Categories `yaml:"categories"`
	Locations      locutil.Locations  `yaml:"locations"`
	PathAndEntries []PathAndEntry     `yaml:"entries"`
}

// PathAndEntry represents how an Entry is displayed in yaml.
type PathAndEntry struct {
	Path   string        `yaml:"path"` // {id}_{surname}-{prename}_{company-name}
	Fields entutil.Entry `yaml:"fields"`
}

// Assemble makes the datastructure that will later be used to create data/entries.yaml.
func Assemble(
	cats catutil.Categories,
	locs locutil.Locations,
	ents entutil.Entries,
) MainListing {

	return MainListing{
		Categories:     cats,
		Locations:      locs,
		PathAndEntries: makePathEntries(ents),
	}
}

func makePathEntries(ents entutil.Entries) []PathAndEntry {
	var paes []PathAndEntry

	for _, ent := range ents {
		for il := range ent.Location { // Generate one per location.
			ent.Title = makeTitle(ent, il)

			c, r, _ := locutil.SplitDisplay(ent.Location[il])
			ent.Country = c
			ent.Region = r

			p := PathAndEntry{
				Path:   makePath(ent, il),
				Fields: ent,
			}
			paes = append(paes, p)
		}
	}

	return paes
}

func makeTitle(f entutil.Entry, locindex int) string {
	var titlePart string
	if (f.Firstname + f.Lastname + f.Credentials) == "" {
		titlePart = f.Company
	} else {
		titlePart = strings.Join([]string{f.Firstname, f.Lastname, f.Credentials}, " ")
	}

	var title string
	countrycode, region, _ := locutil.SplitDisplay(f.Location[locindex])
	if countrycode == "ZZ" || region == "" {
		title = titlePart + fmt.Sprintf(" - %s from %s", f.Category, region)
	} else {
		title = titlePart + fmt.Sprintf(" - %s from %s-%s", f.Category, countrycode, region)
	}

	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "  ", " ")
	return title
}

var regexInvalidPath = regexp.MustCompile("[^A-Za-z0-9_]+")

func makePath(f entutil.Entry, locindex int) string {

	path := fmt.Sprintf("%d_%s-%s_%s",
		f.ID,
		strings.ToLower(f.Firstname),
		strings.ToLower(f.Lastname),
		strings.ToLower(f.Company),
	)
	if locindex > 0 {
		path += fmt.Sprintf("-%d", locindex)
	}

	// Remove diacritics from letters:
	// Cite: https://stackoverflow.com/questions/26722450/remove-diacritics-using-go
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	path, _, _ = transform.String(t, path)

	// Change runs of invalid chars to -
	path = regexInvalidPath.ReplaceAllString(path, "-")
	path = strings.TrimRight(path, "-_") // Clean up the end.

	return path
}

//

// WriteBigYamlFile turns a listing into YAML and writes it to a file in entries.yaml format.
func (m MainListing) WriteBigYamlFile(filename string) error {
	d, err := listingToYaml(m)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, []byte(d), 0666)
}

// listingToYaml turns a listing into the YAML suitable for entries.yaml.
func listingToYaml(listing MainListing) (string, error) {
	d, err := yaml.Marshal(&listing)
	if err != nil {
		return "", err
	}
	dStr := string(d)
	return "---\n" + dStr + "\n", nil
}
