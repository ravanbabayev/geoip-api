FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o geoip-api .

FROM alpine:latest

WORKDIR /root/

COPY .env.example /root/.env

COPY GeoLite2-City.mmdb /app/GeoLite2-City.mmdb

WORKDIR /app

COPY --from=builder /app/geoip-api /geoip-api

EXPOSE 8080

CMD ["/geoip-api"]
