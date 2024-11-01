# Microservices in Go with gRPC , Graphql

## WIP! - Work in progress

Here we are using monorepo kind of Microservices cause all of the Microservices
in this project are identical and they all use pretty much the same packages.
So it's ideal for this use case.
Which might not be ideal for every case of using Microservices.

 Generate Graphql

```sh
cd graphql
go run github.com/99designs/gqlgen generate
```

> Graphql is the entry point of the Microservices

### Top down general control flow

```
From graphql mutation or query ->
client of Microservices ->
server of Microservices ->
service of Microservices ->
repository of Microservices ->
database
```

>Protobuff
>
>- Install protoc
>- Install both plugins protoc-gen-go-grpc and protoc-gen-go
>

```sh
cd account
# go:generate for protoc is defined on account/server.go
go generate # for generating pb files
```

> [!NOTE]
> create pb folder if it's not already there
>

Run Graphql

```sh
cd graphql
go run .
```

## Notes

- Define graphql schema
- Define account model
- Generate graphql with gqlgen

```sh
cd graphql
go run github.com/99designs/gqlgen generate
```

- Define all mutation , query and resolvers in graphql
- setup graphql schema executable using server instance
- popluate config using envconfig in graphql
- use config to serve graphqlserver using NewGraphQLServer
- handle graphql endpoints for graphql and playgorung form gqlgen/handler

- Account define account.proto file
- use proto file and generate pb

```sh
# go:generate for protoc is defined on account/server.go
cd account
go generate # for generating pb files
```

- account microservice built from repository to service to server to client level

- catalog setup
- no db.dockerfile cause will be using elastic search
- [elastic](https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/overview.html)
- should use offical go client for elasticsearch
[elasticsearch](https://github.com/elastic/go-elasticsearch)
[depricated](https://github.com/olivere/elastic/)
