-- name: GetCountry :one
SELECT country_iso2, country_name
FROM countries
WHERE country_iso2 = sqlc.arg(country_iso2);

-- name: GetCodeDetailsByCountryCode :many
SELECT swift_code, address, bank_name, swift_codes.country_iso2
FROM swift_codes
WHERE country_iso2 = sqlc.arg(country_iso2);

-- name: GetCodeDetails :many
SELECT swift_code, address, bank_name, swift_codes.country_iso2, countries.country_name
FROM swift_codes LEFT JOIN countries ON swift_codes.country_iso2 = countries.country_iso2
WHERE swift_codes.swift_code = sqlc.arg(swift_code)
UNION 
SELECT swift_code, address, bank_name, swift_codes.country_iso2, countries.country_name
FROM swift_codes LEFT JOIN countries ON swift_codes.country_iso2 = countries.country_iso2
WHERE RIGHT(sqlc.arg(swift_code), 3) = "XXX"
AND LEFT(swift_code, 8) = LEFT(sqlc.arg(swift_code), 8)
AND NOT RIGHT(swift_code, 3) = "XXX";

-- name: InsertSwiftCode :execresult
INSERT INTO swift_codes (swift_code, address, bank_name, country_iso2)
VALUES (?, ?, ?, ?);

-- name: InsertCountry :execresult
INSERT INTO countries (country_iso2, country_name)
VALUES (?, ?);

-- name: DeleteSwiftCode :execresult
DELETE FROM swift_codes
WHERE swift_code = ?;
