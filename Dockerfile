# syntax=docker/dockerfile:1
# Build
FROM golang:1.18-alpine as build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN  go mod download

COPY . ./

RUN go build -buildvcs=false ./cmd/dyngo


# Deploy
FROM alpine

WORKDIR /

COPY --from=build /app/dyngo /

ENTRYPOINT [ "/dyngo" ]
