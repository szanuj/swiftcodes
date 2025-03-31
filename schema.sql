CREATE TABLE IF NOT EXISTS countries (
    country_iso2 VARCHAR(10) PRIMARY KEY,
    country_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS swift_codes (
    swift_code VARCHAR(50) PRIMARY KEY,
    address TEXT NOT NULL,
    bank_name TEXT NOT NULL,
    country_iso2 VARCHAR(10) NOT NULL,
    FOREIGN KEY (country_iso2) REFERENCES countries (country_iso2)
);
