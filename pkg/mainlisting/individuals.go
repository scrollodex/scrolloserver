package mainlisting

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scrollodex/ResourceUtils/air2hugo/pkg/entutil"
)

const public_entries_dir = "public/entry"
const content_entries_dir = "content/entry"

func (m MainListing) WriteIndividuals() {
	content_entries_dir := "content/entry"
	public_entries_dir := "public/entry"

	// Goroutine: Gather the directory names in public/entry
	var existingPublicDirNames = map[string]struct{}{}
	gatherPublicEntryNames(public_entries_dir, &existingPublicDirNames)
	//fmt.Printf("PUBDIR = %v\n", existingPublicDirNames)

	//Goroutine: In content/entry, gather dirnames and dirname/index.md
	var existingContentEntries = map[string]string{}
	gatherContentEntries(content_entries_dir, &existingContentEntries)
	//fmt.Printf("CONDIR = %v\n", existingContentEntries)

	//func generateIndividualEntries(listingData MainListing, filename string) error {

	for _, pae := range m.PathAndEntries {
		// Generate the filename and content.
		if pae.Fields.Status != 1 {
			// Skip anything marked as not visible.
			//fmt.Printf("DEBUG: skipping status!=1 %q\n", pae.Path)
			continue
		}
		titlePath := pae.Path
		newContents, _ := genJSON(pae.Fields)
		// If content is not in cache or is different, write the contents
		if newContents != existingContentEntries[titlePath] {
			fmt.Printf("UPDATING JSON %s\n", titlePath)
			writeOneEntry(content_entries_dir, titlePath, newContents)
		} // else {
		//fmt.Printf("*** CACHED JSON %s\n", titlePath)
		// }
		// Delete cache from map.
		delete(existingContentEntries, titlePath)
		delete(existingPublicDirNames, titlePath)
	}

	// What's left in the caches is filenames that should be deleted
	// because we haven't generated them (or considered generating them).
	for k := range existingPublicDirNames {
		deletePublic(k)
	}
	for k := range existingContentEntries {
		deleteContent(k)
	}

}

func gatherPublicEntryNames(dirpath string, existingPublicDirNames *(map[string]struct{})) {
	names, err := OSReadDir(dirpath)
	if err != nil {
		log.Fatalf("gatherPublicEntryNames could not ReadDir %q: %s", dirpath, err)
	}
	// If the first char is a digit, this is a Public Entry:
	for _, name := range names {
		if name == "" || name[0] < '0' || name[0] > '9' {
			continue
		}
		(*existingPublicDirNames)[name] = struct{}{}
	}
}

func gatherContentEntries(dirpath string, existingContentEntries *(map[string]string)) {
	names, err := OSReadDir(dirpath)
	if err != nil {
		log.Fatalf("existingContentEntries could not ReadDir %q: %s", dirpath, err)
	}
	// If the first char is a digit, this is a Content Entry:
	for _, name := range names {
		if name == "" || name[0] < '0' || name[0] > '9' {
			continue
		}
		fname := filepath.Join(dirpath, name, "index.md")
		//fmt.Printf("DEBUG: reading %q\n", fname)
		content, err := os.ReadFile(fname)
		if err != nil {
			continue
		}
		(*existingContentEntries)[name] = string(content)
	}
}

func genJSON(data entutil.Entry) (string, error) {
	d, err := json.Marshal(data)
	return string(d), err
}

func writeOneEntry(dirPath, titlePath string, newContents string) error {
	dpath := filepath.Join(dirPath, titlePath)
	err := os.MkdirAll(dpath, os.ModePerm)
	if err != nil {
		return err
	}

	fpath := filepath.Join(dpath, "index.md")
	//fmt.Printf("DEBUG: writing %q with len()=%d\n", fpath, len(newContents))
	err = os.WriteFile(fpath, []byte(newContents), 0o0666)
	//fmt.Printf("DEBUG: write err = %v\n", err)
	return err
}

// namesInDirectory
func OSReadDir(root string) ([]string, error) {
	// Source: https://stackoverflow.com/a/49196644/71978
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func deletePublic(name string) {
	//fmt.Printf("DEBUG: unimplemented delPublic(%q)\n", name)
	fmt.Printf("Deleting DATA %s\n", name)
	dpath := filepath.Join(public_entries_dir, name)
	//fmt.Printf("DEBUG: os.RemoveAll(%q)\n", dpath)
	os.RemoveAll(dpath)
}

func deleteContent(name string) {
	//fmt.Printf("DEBUG: unimplemented delContent(%q)\n", name)
	fmt.Printf("Deleting WEB %s\n", name)
	dpath := filepath.Join(content_entries_dir, name)
	//fmt.Printf("DEBUG: os.RemoveAll(%q)\n", dpath)
	os.RemoveAll(dpath)
}
