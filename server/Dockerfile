# Build Stage
FROM golang:latest AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

CMD ["./bin/server"]
