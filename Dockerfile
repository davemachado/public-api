FROM golang:1.17.6

ENV SRC_DIR=/go/src/github.com/davemachado/public-api

ADD . $SRC_DIR
WORKDIR $SRC_DIR

RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch
COPY --from=0 /go/src/github.com/davemachado/public-api/public-api /public-api
COPY --from=0 /go/src/github.com/davemachado/public-api/static /static
EXPOSE 8080
ENTRYPOINT ["/public-api"]

