package locutil

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/mehanizm/airtable"
)

func Test_convert(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want Location
	}{
		{
			name: "01",
			raw: `
{
	"id": "rec0Y72e0obW1CQIL",
 	"fields": {
		"Created": "2022-02-13T17:55:42.000Z",
		"Created By": {
			"email": "tal@example.org",
			"id": "usrNOCKhVoF3W2Jcv",
			"name": "Tom Limoncelli"
		},
		"Location": "US-OH (Ohio)",
		"x-Comment": "Ohio",
		"x-CountryCode": "US",
		"x-LocationID": 12,
		"x-Region": "OH"
	},
	"createdTime": "2022-02-13T17:55:42.000Z"
}
`,
			want: Location{
				ID:          12,
				Display:     "US-OH (Ohio)",
				CountryCode: "US",
				Region:      "OH",
				Comment:     "Ohio",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var control airtable.Record
			rm := json.RawMessage(tt.raw)
			err := json.Unmarshal(rm, &control)
			if err != nil {
				panic(err)
			}
			if got := convert(&control); !cmp.Equal(got, &(tt.want)) {
				t.Errorf("convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
