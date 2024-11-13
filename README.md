# Microservices in Go with gRPC , Graphql & Elastic Search

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

### [Step by step Notes🔗](/Notes.md)
