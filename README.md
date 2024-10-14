# Microservices in Go with gRPC , Graphql

Here we are using monorepo kind of Microservices cause all of the Microservices
in this project are identical and they all use pretty much the same packages.
So it's ideal for this use case.
Which might not be ideal for every case of using Microservices.

- Generate Graphql

```bash
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
