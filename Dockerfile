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
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-extldflags "-static" -s' -installsuffix cgo -o Mensa-Ingo -v .

# Put everything together
FROM scratch

COPY --from=build-env /go/src/github.com/zottelchin/Mensa-Ingo/Mensa-Ingo /
COPY static /static
COPY 404.html /404.html
COPY mensa.html /mensa.html

ENV GIN_MODE=release
WORKDIR /
EXPOSE 8700

ENTRYPOINT [ "/Mensa-Ingo" ]