package catutil

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
		want Category
	}{
		{
			name: "01",
			raw: `
{
	"id": "recX5DhfmrHwzghF8",
	"fields": {
		"Description": "Physicians, Physicians Assistants, Nurses",
		"IconFilename": "user-md-solid.svg",
		"Name": "Medical Professionals",
		"x-CategoryID": 9,
		"x-Priority": 2
	},
	"createdTime": "2022-02-13T17:55:08.000Z"
}
`,
			want: Category{
				ID:          9,
				Name:        "Medical Professionals",
				Description: "Physicians, Physicians Assistants, Nurses",
				Priority:    2,
				Icon:        "user-md-solid.svg",
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
