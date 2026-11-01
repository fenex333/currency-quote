# currency-quote
Simple exchange rate provider written in Go

## Supported Pairs
- USD/EUR
- USD/MXN
- EUR/USD
- EUR/MXN
- MXN/USD
- MXN/EUR

## Dependencies
- Go >=1.27
- Postgres >=16
- Docker (optional)

## Running

Create ```.env``` file following ```.env.example```, then run commands below
```
# docker run --rm -d -e POSTGRES_PASSWORD=pg -e POSTGRES_DB=pg -p 5432:5432 postgres:16
$ set -a && source .env && set +a
$ go run main.go
```
