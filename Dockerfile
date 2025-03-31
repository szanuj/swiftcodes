FROM golang:1.24.1

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app .

ENV SC_DB_USER="root"
ENV SC_DB_PASSWORD="admin"
ENV SC_DB_NAME="swiftcodes"
ENV SC_DB_HOST="127.0.0.1"
ENV SC_DB_PORT="3306"
ENV SC_API_HOST="127.0.0.1"
ENV SC_API_PORT="8080"

CMD ["app"]