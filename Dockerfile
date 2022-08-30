FROM golang:1.19-bullseye AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./bin/api .

FROM debian:bullseye-slim

LABEL maintainer="Tomash Sidei <tomash.sidey@gmail.com>"

COPY --from=build /app/bin/api /app/api

RUN addgroup --gid 901 spacetrouble && adduser --uid 901 --gid 901 spacetrouble

EXPOSE 8080

USER spacetrouble

ENTRYPOINT ["/app/api"]
