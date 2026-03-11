package dump

import (
	"encoding/json"
	"log"
	"os"
)

// It dumps it into a file for debugging purposes.
func It(filename string, it any) {
	s, err := json.MarshalIndent(it, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(filename, s, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
