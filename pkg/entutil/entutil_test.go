package entutil

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
		want Entry
	}{
		{
			name: "01",
			raw: `
{
	"id": "rec7mcpoY5FCkDOdq",
	"fields": {
		"Address": "MyAddress",
		"Category": "Chiropractors",
		"Company": "MyCompany",
		"Created": "2022-07-12T22:38:38.000Z",
		"Created By": {
			"email": "tal@whatexit.org",
			"id": "usrkunqhVoF3W2Jcv",
			"name": "Tom Limoncelli"
		},
		"Description": "MyDescription",
		"Email": "MyEmail",
		"Email2": "MyEmail2",
		"EntryID": 9999,
		"Fax": "MyFax",
		"Fees": "MyFees",
		"First": "MyFirst",
		"ID": 6001,
		"Job_Title": "MyJobTitle",
		"Last": "MyLast",
		"Last Modified": "2022-07-12T22:40:15.000Z",
		"Last Modified By": {
			"email": "tal@whatexit.org",
			"id": "usrkunqhVoF3W2Jcv",
			"name": "Tom Limoncelli"
		},
		"Location": "AT-All (Austria)",
		"PRIVATE_admin_notes": "MyPrivateAdminNotes",
		"PRIVATE_contact_email": "MyPrivateContactEmail",
		"Phone": "MyPhone",
		"Sal": "MySalutation",
		"Short Description": "MyShortDescription",
		"Status": "SHOW",
		"Suffix": "MySuffix",
		"Website": "MyWebsite",
		"Website2": "MyWebsite2",
		"x-lastUpdate": "2022-01-11",
		"x-private_last_edit_by": "user1"
	},
	"createdTime": "2022-07-12T22:38:38.000Z"
}
`,
			want: Entry{
				ID:        9999,
				Category:  "Chiropractors",
				Locations: []string{"AT-All (Austria)"},
				Status:    1,

				Company:     "MyCompany",
				Salutation:  "MySalutation",
				Firstname:   "MyFirst",
				Lastname:    "MyLast",
				Credentials: "MySuffix",
				JobTitle:    "MyJobTitle",

				ShortDesc:   "MyShortDescription",
				Description: "MyDescription",
				Fees:        "MyFees",

				Address: "MyAddress",
				Email:   "MyEmail",
				Email2:  "MyEmail2",
				Phone:   "MyPhone",
				//Fax:      "MyFax",
				Website:  "MyWebsite",
				Website2: "MyWebsite2",

				LastEditDate: "2022-07-12T22:40:15.000Z",
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
				t.Errorf("convert() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestEntries_Locations(t *testing.T) {
	testdata := Entries{
		{Locations: []string{"repeat"}},
		{Locations: []string{"two"}},
		{Locations: []string{"repeat"}},
	}
	want := []string{"repeat", "two"}

	if got := testdata.Locations(); !cmp.Equal(got, want) {
		t.Errorf("Entries.Locations() = %v, want %v", got, want)
	}
}
