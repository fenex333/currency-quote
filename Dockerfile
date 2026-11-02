# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/currency-quote .

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/currency-quote /usr/local/bin/currency-quote

# edit as needed
ENV EXCHANGE_URL="" \
    EXCHANGE_KEY=""

EXPOSE 8080

ENTRYPOINT ["currency-quote"]
