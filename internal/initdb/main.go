package initdb

import (
	"context"
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strings"

	"swiftcodes/sqlcout"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB_USER      = os.Getenv("SC_DB_USER")
	DB_PASSWORD  = os.Getenv("SC_DB_PASSWORD")
	DB_NAME      = os.Getenv("SC_DB_NAME")
	DB_HOST      = os.Getenv("SC_DB_HOST")
	DB_PORT      = os.Getenv("SC_DB_PORT")
	DB_CONN_BASE = DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/"
)

func ReadCSV(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Couldn't read input file "+path, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	if strings.HasSuffix(path, ".tsv") {
		csvReader.Comma = '\t'
	}
	// csvReader.LazyQuotes = true
	csvReader.FieldsPerRecord = 0
	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Error parsing file "+path, err)
	}
	return rows
}

func ParseData(rows [][]string) ([]sqlcout.InsertSwiftCodeParams, []sqlcout.InsertCountryParams) {
	var codes []sqlcout.InsertSwiftCodeParams
	var countries []sqlcout.InsertCountryParams
	countryMap := make(map[string]string)

	for i := 1; i < len(rows); i++ {
		newCode := sqlcout.InsertSwiftCodeParams{
			SwiftCode:   rows[i][1],
			Address:     rows[i][4],
			BankName:    rows[i][3],
			CountryISO2: strings.ToUpper(rows[i][0]),
		}
		codes = append(codes, newCode)

		_, isKey := countryMap[newCode.CountryISO2]
		if !isKey {
			countryName := strings.ToUpper(rows[i][6])
			countryMap[newCode.CountryISO2] = countryName
			countries = append(countries, sqlcout.InsertCountryParams{
				CountryISO2: newCode.CountryISO2,
				CountryName: countryName,
			})
		}
	}
	return codes, countries
}

func CreateDB(name string, schemaPath string, forTest bool) *sql.DB {
	// Connect to DBMS and create DB
	db, err := sql.Open("mysql", DB_CONN_BASE)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	if forTest {
		db.Exec("DROP DATABASE IF EXISTS " + name)
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		log.Fatal("Failed to create database: ", err)
	}
	db.Close()

	// Connect to our DB specifically
	db, err = sql.Open("mysql", DB_CONN_BASE+name+"?multiStatements=true")
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Read and execute schema
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal("Couldn't read schema file "+schemaPath, err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal("Failed to create tables: ", err)
	}

	return db
}

func PopulateDB(queries *sqlcout.Queries, ctx context.Context, countries []sqlcout.InsertCountryParams, swiftcodes []sqlcout.InsertSwiftCodeParams) {
	for _, country := range countries {
		_, err := queries.InsertCountry(ctx, country)
		if err != nil {
			log.Print("Failed to insert country: ", err)
		}
	}
	for _, code := range swiftcodes {
		_, err := queries.InsertSwiftCode(ctx, code)
		if err != nil {
			log.Print("Failed to insert swift code: ", err)
		}
	}
}

func DBExists(name string) bool {
	db, _ := sql.Open("mysql", DB_CONN_BASE+name)
	err := db.Ping()
	return err == nil
}

func SetupDB(name string, forTest bool) *sql.DB {
	ctx := context.Background()
	db := CreateDB(name, "schema.sql", forTest)

	swiftcodes, countries := ParseData(ReadCSV("swiftcodes.tsv"))

	queries := sqlcout.New(db)

	PopulateDB(queries, ctx, countries, swiftcodes)

	return db
}

func main() {
	SetupDB(DB_NAME, false)
}
