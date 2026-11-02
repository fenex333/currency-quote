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
Create an API key for exchangeratesapi.io <br>
Create ```.env``` file following ```.env.example```, then run commands below
```
$ docker run --rm -d -e POSTGRES_PASSWORD=pg -e POSTGRES_DB=pg -e POSTGRES_USER=pg -p 5432:5432 postgres:16
$ set -a && source .env && set +a
$ go run main.go
```
Alternatively, build and run the Docker image (Postgres is not included and must be run separately)
```
$ docker build -t currency-quote .
$ docker run --rm -d \
    --name currency-quote \
    --network host \
    --env-file .env \
    -p 8080:8080 \
    currency-quote
```

## Endpoint examples
```
/refresh?EUR/USD

/getquote?69c8418b5951604a597d7b70caa5d48602e3eab2

/lastquote?EUR/MXN
```
