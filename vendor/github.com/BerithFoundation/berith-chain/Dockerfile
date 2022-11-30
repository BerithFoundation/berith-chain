# Build Berith in a stock Go builder container
FROM golang:1.13-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /berith
RUN cd /berith && make berith

# Pull Berith into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /berith/build/bin/berith /usr/local/bin/

EXPOSE 8545 8546 8547 40404 40404/udp
ENTRYPOINT ["berith"]
