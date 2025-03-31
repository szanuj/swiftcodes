package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"swiftcodes/internal/initdb"
	"swiftcodes/sqlcout"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	API_VERSION = "1"
	API_NAME    = "swift-codes"
	BASE_URI    = "/v" + API_VERSION + "/" + API_NAME
	COUNTRY     = "country"
)

var (
	DB_USER      = os.Getenv("SC_DB_USER")
	DB_PASSWORD  = os.Getenv("SC_DB_PASSWORD")
	DB_NAME      = os.Getenv("SC_DB_NAME")
	DB_HOST      = os.Getenv("SC_DB_HOST")
	DB_PORT      = os.Getenv("SC_DB_PORT")
	API_HOST     = os.Getenv("SC_API_HOST")
	API_PORT     = os.Getenv("SC_API_PORT")
	DB_CONN_BASE = DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/"
	ctx          context.Context
	db           *sql.DB
	queries      *sqlcout.Queries
)

// Endpoint 1: Retrieve details of a single SWIFT code whether for a headquarters or branches
func GetCodeDetailsHandler(c *gin.Context) {
	swift_code, _ := c.Params.Get("swift_code")
	details, err := queries.GetCodeDetails(ctx, sqlcout.GetCodeDetailsParams{SwiftCode: swift_code})
	if err != nil {
		log.Print("Failed in query GetCodeDetails: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "500 internal server error"})
		return
	}
	if details == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 swift code " + swift_code + " not found"})
		return
	}
	c.JSON(http.StatusOK, MakeDetailsResponse(details))
}

// Endpoint 2: Return all SWIFT codes with details for a specific country (both headquarters and branches)
func GetCodeDetailsByCountryCodeHandler(c *gin.Context) {
	countryISO2, _ := c.Params.Get("country_iso2")
	country, err := queries.GetCountry(ctx, countryISO2)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 country with ISO2 code " + countryISO2 + " not found"})
		return
	}
	details, err := queries.GetCodeDetailsByCountryCode(ctx, countryISO2)
	if err != nil {
		log.Print("Failed in query GetCodeDetailsByCountryCode: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "500 internal server error"})
		return
	}
	c.JSON(http.StatusOK, MakeDetailsByCountryCodeResponse(country, details))
}

// Endpoint 3: Adds new SWIFT code entries to the database for a specific country
func PostSwiftCodeHandler(c *gin.Context) {
	var newCode DetailsInputPayload
	if err := c.BindJSON(&newCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400 bad request structure"})
		return
	}
	if err := ValidateDetailsInputPayload(newCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400 " + err.Error()})
		return
	}
	if _, err := queries.InsertSwiftCode(ctx, sqlcout.InsertSwiftCodeParams{
		Address:     newCode.Address,
		BankName:    newCode.BankName,
		CountryISO2: newCode.CountryISO2,
		SwiftCode:   newCode.SwiftCode,
	}); err != nil {
		msg := strings.Split(err.Error(), " ")
		if len(msg) >= 2 {
			if msg[1] == "1452" {
				c.JSON(http.StatusConflict, gin.H{"error": "409 no country with ISO2 code " + newCode.CountryISO2})
				return
			} else if msg[1] == "1062" {
				c.JSON(http.StatusConflict, gin.H{"error": "409 swift code " + newCode.SwiftCode + " already exists"})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "500 internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "201 swift code " + newCode.SwiftCode + " created"})
}

// Endpoint 4: Deletes swift-code data if swiftCode matches the one in the database
func DeleteSwiftCodeHandler(c *gin.Context) {
	swift_code, _ := c.Params.Get("swift_code")
	result, err := queries.DeleteSwiftCode(ctx, swift_code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "500 internal server error"})
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "500 internal server error"})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 swift code " + swift_code + " not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "200 swift code " + swift_code + " deleted"})
}

func SetupRouter(db_conn_base string, db_name string) (*gin.Engine, error) {
	// Create DB object and check connection
	ctx = context.Background()
	db, err := sql.Open("mysql", db_conn_base+db_name)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	queries = sqlcout.New(db)

	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Link API endpoints
	router.GET(BASE_URI+"/:swift_code", GetCodeDetailsHandler)
	router.GET(BASE_URI+"/country/:country_iso2", GetCodeDetailsByCountryCodeHandler)
	router.POST(BASE_URI, PostSwiftCodeHandler)
	router.DELETE(BASE_URI+"/:swift_code", DeleteSwiftCodeHandler)

	return router, nil
}

func main() {
	if !initdb.DBExists(DB_NAME) {
		db = initdb.SetupDB(DB_NAME, false)
	}

	router, err := SetupRouter(DB_CONN_BASE, DB_NAME)
	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}

	router.Run(API_HOST + ":" + API_PORT)
}
