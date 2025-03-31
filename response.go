package main

import (
	"errors"
	"strings"
	"swiftcodes/sqlcout"
)

type CountryResponse struct {
	CountryISO2 string `json:"countryISO2"`
	CountryName string `json:"countryName"`
}

type DetailsMainResponse struct {
	Address       string                    `json:"address"`
	BankName      string                    `json:"bankName"`
	CountryISO2   string                    `json:"countryISO2"`
	CountryName   string                    `json:"countryName"`
	IsHeadquarter bool                      `json:"isHeadquarter"`
	SwiftCode     string                    `json:"swiftCode"`
	Branches      []DetailsListItemResponse `json:"branches,omitempty"`
}

type DetailsListItemResponse struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

func IsHeadquarter(swiftcode string) bool {
	return strings.HasSuffix(swiftcode, "XXX")
}

func MakeDetailsResponse(details []sqlcout.GetCodeDetailsRow) DetailsMainResponse {
	hq := IsHeadquarter(details[0].SwiftCode)
	response := DetailsMainResponse{
		details[0].Address,
		details[0].BankName,
		details[0].CountryISO2,
		details[0].CountryName.String,
		hq,
		details[0].SwiftCode,
		[]DetailsListItemResponse{},
	}
	for i := 1; i < len(details); i++ {
		response.Branches = append(response.Branches, DetailsListItemResponse{
			details[i].Address,
			details[i].BankName,
			details[i].CountryISO2,
			IsHeadquarter(details[i].SwiftCode),
			details[i].SwiftCode,
		})
	}
	return response
}

type DetailsByCountryCodeResponse struct {
	CountryISO2 string                    `json:"countryISO2"`
	CountryName string                    `json:"countryName"`
	SwiftCodes  []DetailsListItemResponse `json:"swiftCodes"`
}

func MakeDetailsByCountryCodeResponse(country sqlcout.Country, details []sqlcout.SwiftCode) DetailsByCountryCodeResponse {
	response := DetailsByCountryCodeResponse{
		country.CountryISO2,
		country.CountryName,
		[]DetailsListItemResponse{},
	}
	for i := 0; i < len(details); i++ {
		response.SwiftCodes = append(response.SwiftCodes, DetailsListItemResponse{
			details[i].Address,
			details[i].BankName,
			details[i].CountryISO2,
			IsHeadquarter(details[i].SwiftCode),
			details[i].SwiftCode,
		})
	}
	return response
}

type DetailsInputPayload struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

func ValidateDetailsInputPayload(details DetailsInputPayload) error {
	if details.CountryISO2 != strings.ToUpper(details.CountryISO2) {
		return errors.New("countryISO2 must be uppercase")
	}
	if details.CountryName != strings.ToUpper(details.CountryName) {
		return errors.New("countryName must be uppercase")
	}
	if details.IsHeadquarter != IsHeadquarter(details.SwiftCode) {
		return errors.New("swiftCode must end with 'XXX' if and only if it is headquarter")
	}
	return nil
}
