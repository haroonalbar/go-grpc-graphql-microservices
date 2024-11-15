FROM golang:1.23.2-alpine3.20 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/haroonalbar/go-grpc-graphql-microservices
COPY go.mod go.sum ./
COPY vendor vendor
COPY catalog catalog 
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./catalog/cmd/catalog

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD [ "app" ]
