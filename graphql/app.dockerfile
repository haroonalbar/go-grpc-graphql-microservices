FROM golang:1.23.2-alpine3.20 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/haroonalbar/go-grpc-graphql-microservices
COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
COPY catalog catalog
COPY order order
COPY graphql graphql
COPY static static
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./graphql

FROM alpine:3.20
WORKDIR /usr/bin
COPY static static
COPY --from=build /go/bin .
EXPOSE 8080
CMD [ "app" ]