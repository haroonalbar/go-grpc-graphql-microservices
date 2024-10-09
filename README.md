# Microservices in Go with gRPC , Graphql

- Generate Graphql

```bash
cd graphql
go run github.com/99designs/gqlgen generate
```

> Graphql is the entry point of the Microservices

### Top level general understanding

```
From graphql mutation or query ->
client of Microservices ->
server of Microservices ->
service of Microservices ->
repository of Microservices ->
database
```
