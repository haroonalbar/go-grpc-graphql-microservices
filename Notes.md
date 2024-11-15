## Notes for myself

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
- [Elastic Search](https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/overview.html)
- should use offical go client for elasticsearch
[elasticsearch](https://github.com/elastic/go-elasticsearch)
[depricated](https://github.com/olivere/elastic/)
- setup all things repository, service , proto, server
- for pb genaration

```sh
cd catalog
go generate #it's defined on catalog/server
```

- Order:
- setup sql.up , db, app dockerfile
- setup Repository
- setup Service
- setup order.proto
- add generate comment to service.go
- generate pb

```sh
cd order
go generate
```

- setup main - connect to db and listen and serve order microservice on 8080
- setup Server and client

- docker compose file done

- graphql:
- update schema
- setup generate directive for gqlgen

```sh
cd graphql
go generate # located in ./graphql/graph.go
```

- envconfig not working perfectly

```bash
# First, clear any existing env vars     
unset ACCOUNT_SERVICE_URL CATALOG_SERVICE_URL ORDER_SERVICE_URL

# Then set them one by one
export ACCOUNT_SERVICE_URL="http://localhost:8081"
export CATALOG_SERVICE_URL="http://localhost:8082"
export ORDER_SERVICE_URL="http://localhost:8083"

# Verify they're set
env | grep URL

# Run the application
go run .
```

- Go to grahql playground
<http://localhost:8080/playground>

- Docker:
- Make sure docker daemon is running
- In my case OrbStack (you can also use Docker desktop)

- compose docker

```bash
docker-compse up -d --build
```

- Lot of things went wrong
- had to debug and fix a lot of things
