package dump

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// It dumps it into a file for debugging purposes.
func It(filename string, it interface{}) {
	s, err := json.MarshalIndent(it, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filename, s, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
