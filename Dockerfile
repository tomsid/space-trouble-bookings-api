FROM golang:1.19-bullseye AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/api .

FROM debian:bullseye-slim

LABEL maintainer="Tomash Sidei <tomash.sidey@gmail.com>"

COPY --from=build /app/bin/api /app/api
COPY db_migrations /db_migrations

RUN addgroup --gid 901 spacetrouble && adduser --uid 901 --gid 901 spacetrouble

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 8080

USER spacetrouble

ENTRYPOINT ["/app/api"]
