FROM golang:1.12 as base

MAINTAINER "Hunter Long (https://github.com/hunterlong)"

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY *.go ./

RUN go build -a -tags netgo -ldflags "-w -extldflags '-static'" -o /app/gethexporter

FROM alpine:latest

MAINTAINER "Hunter Long (https://github.com/hunterlong)"
RUN apk add --no-cache ca-certificates
COPY --from=base /app/gethexporter /usr/local/bin/gethexporter

ENV GETH https://mainnet.infura.io/v3/f5951d9239964e62aa32ca40bad376a6
ENV ADDRESSES ""
ENV DELAY 500

EXPOSE 9090
ENTRYPOINT gethexporter
