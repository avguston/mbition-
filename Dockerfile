##
## Build
##
FROM golang:1.16-buster AS build

WORKDIR /app

COPY . /app

RUN go mod download

COPY *.go ./

RUN go test -v -coverprofile=/dev/stdout -coverpkg=github.com/apiresponse . && \
    go build -o /buser

##
## Deploy
##
FROM registry.access.redhat.com/ubi7/ubi-minimal:7.9

WORKDIR /

COPY --from=build /buser /buser
COPY --from=build /app/templates /templates

EXPOSE 8080

ENTRYPOINT ["/buser"]
