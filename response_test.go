package main

import (
	"database/sql"
	"reflect"
	"swiftcodes/sqlcout"
	"testing"
)

func TestIsHeadquarter(t *testing.T) {
	tt := []struct {
		swiftcode string
		want      bool
	}{
		{"", false},
		{"a", false},
		{"XXX", true},
		{"BCRCBGS1XXX", true},
		{"BJSBMCMXLCO", false},
		{"XXXBMC MXXXXA", false},
	}
	for i := 0; i < len(tt); i++ {
		out := IsHeadquarter(tt[i].swiftcode)
		if out != tt[i].want {
			t.Errorf(`IsHeadquarter("%s") = %t, want %t`, tt[i].swiftcode, out, tt[i].want)
		}
	}
}

func TestMakeDetailsResponse(t *testing.T) {
	tt := []struct {
		details []sqlcout.GetCodeDetailsRow
		want    DetailsMainResponse
	}{
		{
			[]sqlcout.GetCodeDetailsRow{
				{
					SwiftCode:   "A",
					Address:     "A",
					BankName:    "A",
					CountryISO2: "A",
					CountryName: sql.NullString{String: "A", Valid: true},
				},
			},
			DetailsMainResponse{
				SwiftCode:     "A",
				Address:       "A",
				BankName:      "A",
				CountryISO2:   "A",
				CountryName:   "A",
				IsHeadquarter: false,
				Branches:      []DetailsListItemResponse{},
			},
		},
		{
			[]sqlcout.GetCodeDetailsRow{
				{
					SwiftCode:   "AAAAAAAAXXX",
					Address:     "A",
					BankName:    "A",
					CountryISO2: "A",
					CountryName: sql.NullString{String: "A", Valid: true},
				},
				{
					SwiftCode:   "AAAAAAAABBB",
					Address:     "A",
					BankName:    "A",
					CountryISO2: "A",
					CountryName: sql.NullString{String: "A", Valid: true},
				},
			},
			DetailsMainResponse{
				SwiftCode:     "AAAAAAAAXXX",
				Address:       "A",
				BankName:      "A",
				CountryISO2:   "A",
				CountryName:   "A",
				IsHeadquarter: true,
				Branches: []DetailsListItemResponse{
					{
						SwiftCode:     "AAAAAAAABBB",
						Address:       "A",
						BankName:      "A",
						CountryISO2:   "A",
						IsHeadquarter: false,
					},
				},
			},
		},
	}
	for i := 0; i < len(tt); i++ {
		out := MakeDetailsResponse(tt[i].details)
		if !reflect.DeepEqual(out, tt[i].want) {
			t.Errorf(`MakeDetailsResponse("%v") = %v, want %v`, tt[i].details, out, tt[i].want)
		}
	}
}

func TestMakeDetailsByCountryCodeResponse(t *testing.T) {
	tt := []struct {
		country   sqlcout.Country
		swiftcode []sqlcout.SwiftCode
		want      DetailsByCountryCodeResponse
	}{
		{
			sqlcout.Country{
				CountryISO2: "WT",
				CountryName: "WATANIA",
			},
			[]sqlcout.SwiftCode{
				{
					SwiftCode:   "AXXX",
					Address:     "",
					BankName:    "",
					CountryISO2: "WT",
				},
				{
					SwiftCode:   "B",
					Address:     "",
					BankName:    "",
					CountryISO2: "WT",
				},
			},
			DetailsByCountryCodeResponse{
				CountryISO2: "WT",
				CountryName: "WATANIA",
				SwiftCodes: []DetailsListItemResponse{
					{
						SwiftCode:     "AXXX",
						Address:       "",
						BankName:      "",
						CountryISO2:   "WT",
						IsHeadquarter: true,
					},
					{
						SwiftCode:     "B",
						Address:       "",
						BankName:      "",
						CountryISO2:   "WT",
						IsHeadquarter: false,
					},
				},
			},
		},
		{
			sqlcout.Country{
				CountryISO2: "WT",
				CountryName: "WATANIA",
			},
			[]sqlcout.SwiftCode{},
			DetailsByCountryCodeResponse{
				CountryISO2: "WT",
				CountryName: "WATANIA",
				SwiftCodes:  []DetailsListItemResponse{},
			},
		},
	}
	for i := 0; i < len(tt); i++ {
		out := MakeDetailsByCountryCodeResponse(tt[i].country, tt[i].swiftcode)
		if !reflect.DeepEqual(out, tt[i].want) {
			t.Errorf(`MakeDetailsByCountryCodeResponse("%v", "%v") = %v, want %v`, tt[i].country, tt[i].swiftcode, out, tt[i].want)
		}
	}
}

func TestValidateDetailsInputPayload(t *testing.T) {
	tt := []struct {
		details DetailsInputPayload
		wantErr bool
	}{
		{
			DetailsInputPayload{
				Address:       "",
				BankName:      "",
				CountryISO2:   "",
				CountryName:   "",
				IsHeadquarter: true,
				SwiftCode:     "",
			},
			true,
		},
		{
			DetailsInputPayload{
				Address:       "",
				BankName:      "",
				CountryISO2:   "",
				CountryName:   "",
				IsHeadquarter: true,
				SwiftCode:     "XXX",
			},
			false,
		},
		{
			DetailsInputPayload{
				Address:       "",
				BankName:      "",
				CountryISO2:   "wt",
				CountryName:   "",
				IsHeadquarter: false,
				SwiftCode:     "",
			},
			true,
		},
		{
			DetailsInputPayload{
				Address:       "",
				BankName:      "",
				CountryISO2:   "WT",
				CountryName:   "",
				IsHeadquarter: false,
				SwiftCode:     "",
			},
			false,
		},
		{
			DetailsInputPayload{
				Address:       "",
				BankName:      "",
				CountryISO2:   "",
				CountryName:   "Watania4",
				IsHeadquarter: false,
				SwiftCode:     "",
			},
			true,
		},
		{
			DetailsInputPayload{
				Address:       "",
				BankName:      "",
				CountryISO2:   "",
				CountryName:   "WATANIA4",
				IsHeadquarter: false,
				SwiftCode:     "",
			},
			false,
		},
		{
			DetailsInputPayload{
				Address:       "",
				BankName:      "",
				CountryISO2:   "",
				CountryName:   "WATANIA4",
				IsHeadquarter: false,
				SwiftCode:     "AAAXXXX",
			},
			true,
		},
	}
	for i := 0; i < len(tt); i++ {
		err := ValidateDetailsInputPayload(tt[i].details)
		if err == nil && tt[i].wantErr {
			t.Errorf(`ValidateDetailsInputPayload("%v") = nil, wanted error`, tt[i].details)
		} else if err != nil && !tt[i].wantErr {
			t.Errorf(`ValidateDetailsInputPayload("%v") = error %v, wanted nil`, tt[i].details, err)
		}
	}
}
