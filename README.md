# Microservices in Go with gRPC, GraphQL & Elasticsearch

## Overview

This project demonstrates a microservices architecture utilizing gRPC for inter-service communication and GraphQL as the API gateway. It includes services for account management, product catalog, and order processing. The architecture is designed to be modular and scalable, making it suitable for various applications.

## Technologies Used

- [**Go**](https://golang.org/): The primary programming language for building the microservices.
- [**gRPC**](https://grpc.io/docs/): A high-performance RPC framework for communication between services.
- [**GraphQL**](https://graphql.org/learn/): A query language for APIs, providing a more efficient and flexible alternative to REST.
- [**Elasticsearch**](https://www.elastic.co/): A distributed search and analytics engine used for the catalog service.
- [**PostgreSQL**](https://www.postgresql.org/): A relational database used for account and order services.
- [**Docker**](https://www.docker.com/): For containerization and orchestration of services.

## Project Structure

The project consists of the following main components:

- **Account Service**: Manages user accounts and authentication.
- **Catalog Service**: Handles product listings and search functionalities.
- **Order Service**: Manages order processing and transactions.
- **GraphQL API Gateway**: Acts as the entry point for client requests, routing them to the appropriate services.

Each service has its own database:

- Account and Order services use PostgreSQL.
- Catalog service uses Elasticsearch.

All services file structure are similar.

```
.
├── account.proto
├── app.dockerfile
├── client.go
├── cmd
│   └── account
│       └── main.go
├── db.dockerfile
├── pb
│   ├── account.pb.go
│   └── account_grpc.pb.go
├── repository.go
├── server.go
├── service.go
└── up.sql
```

## Getting Started

1. **Clone the repository**:

   ```bash
   git clone https://github.com/haroonalbar/go-grpc-graphql-microservices
   cd go-grpc-graphql-microservices
   ```

2. **Start the services using Docker Compose**:

   ```bash
   docker-compose up -d --build
   ```

3. **Access the GraphQL playground** at `http://localhost:8080/playground`.

 > **Or access the demo htmx frontend** at `http://localhost:8080`.

## GraphQL API Usage

The GraphQL API provides a unified interface to interact with all the microservices.

### Example Queries and Mutations

#### Query Accounts

```graphql
query {
  accounts {
    id
    name
  }
}
```

#### Create an Account

```graphql
mutation {
  createAccount(account: {name: "New Account"}) {
    id
    name
  }
}
```

#### Query Products

```graphql
query {
  products {
    id
    name
    price
  }
}
```

#### Create a Product

```graphql
mutation {
  createProduct(product: {name: "New Product", description: "A new product", price: 19.99}) {
    id
    name
    price
  }
}
```

#### Create an Order

```graphql
mutation {
  createOrder(order: {accountId: "account_id", products: [{id: "product_id", quantity: 2}]}) {
    id
    totalPrice
    products {
      name
      quantity
    }
  }
}
```

### Advanced Queries

#### Pagination and Filtering

```graphql
query {
  products(pagination: {skip: 0, take: 5}, query: "search_term") {
    id
    name
    description
    price
  }
}
```

#### Calculate Total Spent by an Account

```graphql
query {
  accounts(id: "account_id") {
    name
    orders {
      totalPrice
    }
  }
}
```

## gRPC File Generation

To generate gRPC files, follow these steps:

1. Download and install [protoc](https://grpc.io/docs/protoc-installation)

2. Install the necessary Go plugins:

   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. Create the `pb` folder in your service if doesn't exist.

   ```bash
    mkdir pb
   ```

5. Finally, run the command to generate the files:

   ```bash
   # The Go generate directive is defined in the on graph.go on graphql service
   # and on service.go of all the other services
   go generate
   ```

## Acknowledgments

Special thanks to [@AkhilSharma90](https://github.com/AkhilSharma90) for the valuable insights and resources that contributed to the development of this project.

## Conclusion

This project serves as a comprehensive example of building a microservices architecture using Go, gRPC, GraphQL, and Elasticsearch. It provides a solid foundation for further development and scaling of microservices-based applications.
