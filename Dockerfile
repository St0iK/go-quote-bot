# build stage
FROM golang:1.11.4-alpine3.8 AS build-env

WORKDIR /go/src/github.com/St0iK/go-quote-bot

COPY ./ .

RUN apk --no-cache add git bzr mercurial && \
    go get -u github.com/golang/dep/... && \
    dep ensure -v --vendor-only && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o quote-bot .

# -------------------------------------------------------------------------------
# final stage
FROM alpine:latest  

ARG MONGO_DB_URL
ENV MONGO_DB_URL ${MONGO_DB_URL}

WORKDIR /root/

COPY --from=build-env /go/src/github.com/St0iK/go-quote-bot .

RUN apk --no-cache add ca-certificates

# ENTRYPOINT tail -f /dev/null
CMD ["./quote-bot"]
