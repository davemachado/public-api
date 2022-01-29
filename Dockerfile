FROM golang:1.17.6

ENV SRC_DIR=/go/src/github.com/davemachado/public-api

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && chmod +x /usr/local/bin/dep

ADD . $SRC_DIR
WORKDIR $SRC_DIR

RUN go build

EXPOSE 8080
ENTRYPOINT ["./public-api"]
