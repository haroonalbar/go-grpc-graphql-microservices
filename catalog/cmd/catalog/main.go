package main

import (
	"log"

	"github.com/haroonalbar/go-grpc-graphql-microservices/catalog"
)

func main() {
	r, err := catalog.NewElasticRepository("")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
}
