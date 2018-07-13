FROM golang

ADD . /go/src/github.com/hunterlong/gethexporter
RUN cd /go/src/github.com/hunterlong/gethexporter && go get
RUN go install github.com/hunterlong/gethexporter

ENV GETH https://mainnet.infura.io/3dek1oWkaoYCx6uXRGYv
ENV PORT 9090

RUN mkdir /app
WORKDIR /app

EXPOSE 9090

ENTRYPOINT /go/bin/gethexporter
