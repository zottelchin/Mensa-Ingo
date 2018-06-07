# This Dockerfile builds everything at once.
# It can be used for an automated build on the Docker Hub.

# Build
FROM golang:latest AS build-env
RUN go get -u github.com/golang/dep/cmd/dep
RUN mkdir -p /go/src/github.com/zottelchin/Mensa-Ingo
WORKDIR /go/src/github.com/zottelchin/Mensa-Ingo
COPY *.go /go/src/github.com/zottelchin/Mensa-Ingo/
COPY Gopkg.* /go/src/github.com/zottelchin/Mensa-Ingo/
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static" -s' -installsuffix cgo -o Mensa-Ingo -v .

FROM alpine:latest as network
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
RUN zip -r -0 /zoneinfo.zip .

# Put everything together
FROM scratch

COPY --from=build-env /go/src/github.com/zottelchin/Mensa-Ingo/Mensa-Ingo /
COPY static /static
COPY 404.html /404.html
COPY mensa.html /mensa.html
COPY --from=network /zoneinfo.zip /
COPY --from=network /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV GIN_MODE=release
WORKDIR /
EXPOSE 8700

ENTRYPOINT [ "/Mensa-Ingo" ]