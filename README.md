# currency-quote
Simple currency quote server written in Go

## Dependencies
- Go >=1.27
- Postgres >=16
- Docker (optional)

## Running

Create .env file following .env.example, then run commands below
```
$ docker run --rm -d --name postgres -p 5432:5432 postgres:16
$ source .env
$ go run main.go
```
