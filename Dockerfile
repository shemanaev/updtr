# syntax=docker/dockerfile:1

## Build
FROM golang:1.19.1-bullseye AS build
ARG VERSION=dev
ARG BUILD_DATE=dev
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X github.com/shemanaev/updtr/internal/meta.Version=$VERSION -X github.com/shemanaev/updtr/internal/meta.date=$BUILD_DATE" -o /go/bin/updtr

## Deploy
FROM gcr.io/distroless/base-debian11

COPY --from=build /go/bin/updtr /bin/updtr

EXPOSE 8080/tcp
CMD ["updtr"]
