package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/kelseyhightower/envconfig"
)

// The envconfig tags you see in this Go struct definition are part of a popular
// Go library called "envconfig" (or sometimes "go-envconfig").
// This library is used to populate struct fields from environment variables.
// Here's what these tags do:
// - They map environment variables to struct fields.
// - They allow you to specify the names of the environment variables that should be used to populate each field.
type AppConfig struct {
	AccountUrl string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogUrl string `envconfig:"CATALOG_SERVICE_URL"`
	OrderUrl   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	// Process populates the specified struct based on environment variables
	// TODO:
	// - envconfig no longer mainted change the implementation
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Error populating env : %v", err)
	}

	// graph.go
	// Create a Graphql server
	s, err := NewGraphQLServer(cfg.AccountUrl, cfg.CatalogUrl, cfg.OrderUrl)
	if err != nil {
		log.Fatalf("Error setting Graphql server: %v", err)
	}

	// NOTE:
	// updated to new serve mux instead of default one
	mux := http.NewServeMux()

	// sets up the main GraphQL endpoint where clients can send queries and mutations.
	mux.Handle("/graphql", handler.New(s.ToExecutableSchema()))
	// provides a web-based GraphQL Playground interface for easy testing and exploration of the GraphQL API.
	mux.Handle("/playground", playground.Handler("play", "/graphql"))

	// Run server
	log.Fatal(http.ListenAndServe(":8080", mux))
}
