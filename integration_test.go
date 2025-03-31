package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"swiftcodes/internal/initdb"
)

const TEST_DB_NAME = "test"

func JSONEqual(a, b string) (bool, error) {
	var aInterface, bInterface interface{}
	if err := json.Unmarshal([]byte(a), &aInterface); err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(b), &bInterface); err != nil {
		return false, err
	}
	return reflect.DeepEqual(aInterface, bInterface), nil
}

func TestGetCodeDetailsHandler(t *testing.T) {
	db := initdb.SetupDB(TEST_DB_NAME, true)
	defer db.Exec("DROP DATABASE IF EXISTS " + TEST_DB_NAME)
	router, err := SetupRouter(DB_CONN_BASE, TEST_DB_NAME)
	if err != nil {
		t.Errorf("TestGetCodeDetailsHandler() DB connection error: %v", err)
	}

	tt := []struct {
		method       string
		url          string
		reader       io.Reader
		wantCode     int
		wantResponse string
	}{
		{
			http.MethodGet,
			"/v1/swift-codes/ABC",
			nil,
			http.StatusNotFound,
			`{"error":"404 swift code ABC not found"}`,
		},
		{
			http.MethodGet,
			"/v1/swift-codes/BIGBPLPWCUS",
			nil,
			http.StatusOK,
			`{"address":"HARMONY CENTER UL. STANISLAWA ZARYNA 2A WARSZAWA, MAZOWIECKIE, 02-593","bankName":"BANK MILLENNIUM S.A.","countryISO2":"PL","countryName":"POLAND","isHeadquarter":false,"swiftCode":"BIGBPLPWCUS"}`,
		},
		{
			http.MethodGet,
			"/v1/swift-codes/BIGBPLPWXXX",
			nil,
			http.StatusOK,
			`{"address":"HARMONY CENTER UL. STANISLAWA ZARYNA 2A WARSZAWA, MAZOWIECKIE, 02-593","bankName":"BANK MILLENNIUM S.A.","countryISO2":"PL","countryName":"POLAND","isHeadquarter":true,"swiftCode":"BIGBPLPWXXX","branches":[{"address":"HARMONY CENTER UL. STANISLAWA ZARYNA 2A WARSZAWA, MAZOWIECKIE, 02-593","bankName":"BANK MILLENNIUM S.A.","countryISO2":"PL","isHeadquarter":false,"swiftCode":"BIGBPLPWCUS"}]}`,
		},
	}

	for i := 0; i < len(tt); i++ {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(tt[i].method, tt[i].url, tt[i].reader)
		if err != nil {
			t.Errorf("TestGetCodeDetailsHandler() error handling request: %v", err)
		}
		router.ServeHTTP(w, req)

		if w.Code != tt[i].wantCode {
			t.Errorf("TestGetCodeDetailsHandler() response code %v, want %v",
				w.Code, tt[i].wantCode)
		}
		responseCorrect, err := JSONEqual(tt[i].wantResponse, w.Body.String())
		if err != nil {
			t.Errorf("TestGetCodeDetailsHandler() error comparing JSON actual response %v wanted response %v: ",
				w.Body.String(), tt[i].wantResponse)
		}
		if !responseCorrect {
			t.Errorf("TestGetCodeDetailsHandler() response code %v, want %v",
				w.Body.String(), tt[i].wantResponse)
		}
	}
}

func TestGetCodeDetailsByCountryCodeHandler(t *testing.T) {
	db := initdb.SetupDB(TEST_DB_NAME, true)
	defer db.Exec("DROP DATABASE IF EXISTS " + TEST_DB_NAME)
	router, err := SetupRouter(DB_CONN_BASE, TEST_DB_NAME)
	if err != nil {
		t.Errorf("TestGetCodeDetailsByCountryCodeHandler() DB connection error: %v", err)
	}

	tt := []struct {
		method   string
		url      string
		reader   io.Reader
		wantCode int
	}{
		{
			http.MethodGet,
			"/v1/swift-codes/country/abc",
			nil,
			http.StatusNotFound,
		},
		{
			http.MethodGet,
			"/v1/swift-codes/country/mt",
			nil,
			http.StatusOK,
		},
		{
			http.MethodGet,
			"/v1/swift-codes/country/pl",
			nil,
			http.StatusOK,
		},
	}

	for i := 0; i < len(tt); i++ {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(tt[i].method, tt[i].url, tt[i].reader)
		if err != nil {
			t.Errorf("TestGetCodeDetailsByCountryCodeHandler() error handling request: %v", err)
		}
		router.ServeHTTP(w, req)

		if w.Code != tt[i].wantCode {
			t.Errorf("TestGetCodeDetailsByCountryCodeHandler() response code %v, want %v",
				w.Code, tt[i].wantCode)
		}
	}
}

func TestPostSwiftCodeHandler(t *testing.T) {
	db := initdb.SetupDB(TEST_DB_NAME, true)
	defer db.Exec("DROP DATABASE IF EXISTS " + TEST_DB_NAME)
	router, err := SetupRouter(DB_CONN_BASE, TEST_DB_NAME)
	if err != nil {
		t.Errorf("TestGetCodeDetailsByCountryCodeHandler() DB connection error: %v", err)
	}

	tt := []struct {
		method   string
		url      string
		payload  DetailsInputPayload
		wantCode int
	}{
		{
			http.MethodPost,
			"/v1/swift-codes",
			DetailsInputPayload{
				Address:       "A",
				BankName:      "A",
				CountryISO2:   "AL",
				CountryName:   "ALBANIA",
				IsHeadquarter: false,
				SwiftCode:     "AAA",
			},
			http.StatusCreated,
		},
		{
			http.MethodPost,
			"/v1/swift-codes",
			DetailsInputPayload{
				Address:       "A",
				BankName:      "A",
				CountryISO2:   "AL",
				CountryName:   "ALBANIA",
				IsHeadquarter: false,
				SwiftCode:     "AAA",
			},
			http.StatusConflict,
		},
		{
			http.MethodPost,
			"/v1/swift-codes",
			DetailsInputPayload{
				Address:       "A",
				BankName:      "A",
				CountryISO2:   "ABC",
				CountryName:   "ABC",
				IsHeadquarter: false,
				SwiftCode:     "BBB",
			},
			http.StatusConflict,
		},
		{
			http.MethodPost,
			"/v1/swift-codes",
			DetailsInputPayload{
				Address:       "A",
				BankName:      "A",
				CountryISO2:   "AL",
				CountryName:   "ALBANIA",
				IsHeadquarter: false,
				SwiftCode:     "XXX",
			},
			http.StatusBadRequest,
		},
		{
			http.MethodPost,
			"/v1/swift-codes",
			DetailsInputPayload{
				Address:       "A",
				BankName:      "A",
				CountryISO2:   "al",
				CountryName:   "albania",
				IsHeadquarter: false,
				SwiftCode:     "ABC",
			},
			http.StatusBadRequest,
		},
	}

	for i := 0; i < len(tt); i++ {
		w := httptest.NewRecorder()
		jsonPayload, _ := json.Marshal(tt[i].payload)
		req, err := http.NewRequest(tt[i].method, tt[i].url, strings.NewReader(string(jsonPayload)))
		if err != nil {
			t.Errorf("TestPostSwiftCodeHandler() error handling request: %v", err)
		}
		router.ServeHTTP(w, req)

		if w.Code != tt[i].wantCode {
			t.Errorf("TestPostSwiftCodeHandler() test index %v. response code %v, want %v",
				i, w.Code, tt[i].wantCode)
		}
	}
}

func TestDeleteSwiftCodeHandler(t *testing.T) {
	db := initdb.SetupDB(TEST_DB_NAME, true)
	defer db.Exec("DROP DATABASE IF EXISTS " + TEST_DB_NAME)
	router, err := SetupRouter(DB_CONN_BASE, TEST_DB_NAME)
	if err != nil {
		t.Errorf("TestDeleteSwiftCodeHandler() DB connection error: %v", err)
	}

	tt := []struct {
		method   string
		url      string
		reader   io.Reader
		wantCode int
	}{
		{
			http.MethodDelete,
			"/v1/swift-codes/TEST",
			nil,
			http.StatusNotFound,
		},
		{
			http.MethodDelete,
			"/v1/swift-codes/BIGBPLPWCUS",
			nil,
			http.StatusOK,
		},
	}

	for i := 0; i < len(tt); i++ {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(tt[i].method, tt[i].url, tt[i].reader)
		if err != nil {
			t.Errorf("TestDeleteSwiftCodeHandler() error handling request: %v", err)
		}
		router.ServeHTTP(w, req)

		if w.Code != tt[i].wantCode {
			t.Errorf("TestDeleteSwiftCodeHandler() test index %v. response code %v, want %v",
				i, w.Code, tt[i].wantCode)
		}
	}
}
