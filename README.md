# Swift Codes Exercise

## Setup instructions

### Building

- Clone this repository
- Download dependencies using `go mod download`
- Execute `go build .` in project directory

### Running

A locally running instance of MariaDB on port `3306` with root password `admin` is required for this app to run. You can easily set it up via Docker:

- `docker pull mariadb`, then
- `docker run --name mariadbtest -e MYSQL_ROOT_PASSWORD=admin -p 3306:3306 -d`

You will need to run the app with environment variables from the .env file. Example using godotenv:

- First install godotenv `go install github.com/joho/godotenv/cmd/godotenv@latest`
- Run app via `godotenv -f go run .`

### Testing

Run unit and integration tests via

- `godotenv -f .env go test .`
