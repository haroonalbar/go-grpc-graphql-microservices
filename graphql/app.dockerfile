# Use the official Golang image from the Docker Hub, version 1.23.2 based on Alpine Linux 3.20 as the build stage
FROM golang:1.23.2-alpine3.20 AS build
# Install necessary packages like GCC, G++, make, and certificates
RUN apk --no-cache add gcc g++ make ca-certificates
# Set the working directory inside the container
WORKDIR /go/src/github.com/haroonalbar/go-grpc-graphql-microservices
# Copy the Go module and sum files
COPY go.mod go.sum ./
# Copy the vendor folder which contains all dependencies
COPY vendor vendor
# Copy the account service directory
COPY account account
# Copy the catalog service directory
COPY catalog catalog
# Copy the order service directory
COPY order order
# Copy the graphql service directory
COPY graphql graphql
# Compile the application to a binary named app using the vendor folder for dependencies
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./graphql


# Use a smaller image for the final stage to reduce size
FROM alpine:3.11
# Set the working directory to /usr/bin
WORKDIR /usr/bin
# Copy the compiled binary from the build stage
COPY --from=build /go/bin .
# Expose port 8080 for the application
EXPOSE 8080
# Command to run the binary
CMD [ "app" ]
