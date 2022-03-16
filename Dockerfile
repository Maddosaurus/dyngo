# syntax=docker/dockerfile:1
# See https://blog.baeke.info/2021/03/28/distroless-or-scratch-for-go-apps/

ARG GO_VERSION=1.18

#
# Build
#
FROM golang:${GO_VERSION}-alpine as build

RUN apk add --no-cache git

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN  go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build -buildvcs=false -installsuffix 'static' -o /dyngo ./cmd/dyngo


#
# Deploy
#
FROM gcr.io/distroless/static AS final

LABEL maintainer="Maddosaurus"
USER nonroot:nonroot

COPY --from=build --chown=nonroot:nonroot /dyngo /dyngo

ENTRYPOINT [ "/dyngo" ]
