FROM golang:1.8

ENV SRC_DIR=/go/src/github.com/davemachado/public-api

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && chmod +x /usr/local/bin/dep

ADD . $SRC_DIR
WORKDIR $SRC_DIR

RUN dep ensure -vendor-only
RUN curl -fsSl -o entries.min.json https://raw.githubusercontent.com/toddmotto/public-apis/master/json/entries.min.json 

RUN go build

EXPOSE 8080
ENTRYPOINT ["./public-api"]
