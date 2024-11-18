package main

import (
	"log"
	"net/http"

	// "github.com/99designs/gqlgen/handler"

	"github.com/99designs/gqlgen/handler"
	"github.com/kelseyhightower/envconfig"
)

// WARN: on schema.graphql products method of tye Query the extra added may or maynot
// cause issues later

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
	// log.Printf("Environment variables as seen by OS:")
	// log.Printf("ACCOUNT_URL: %s", os.Getenv("ACCOUNT_SERVICE_URL"))
	// log.Printf("CATALOG_URL: %s", os.Getenv("CATALOG_SERVICE_URL"))
	// log.Printf("ORDER_URL: %s", os.Getenv("ORDER_SERVICE_URL"))

	var cfg AppConfig
	// Process populates the specified struct based on environment variables
	// TODO: envconfig no longer mainted change the implementation
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Error populating env : %v", err)
	}

	log.Printf("Config after processing:")
	log.Printf("AccountUrl: %s", cfg.AccountUrl)
	log.Printf("CatalogUrl: %s", cfg.CatalogUrl)
	log.Printf("OrderUrl: %s", cfg.OrderUrl)

	if cfg.AccountUrl == "" {
		log.Println("AccountUrl is empty in configuration")
	}
	if cfg.CatalogUrl == "" {
		log.Println("CatalogUrl is empty in configuration")
	}
	if cfg.OrderUrl == "" {
		log.Println("OrderUrl is empty in configuration")
	}

	// graph.go
	// Create a Graphql server
	s, err := NewGraphQLServer(cfg.AccountUrl, cfg.CatalogUrl, cfg.OrderUrl)
	if err != nil {
		log.Fatalf("Error setting Graphql server: %v", err)
	}

	// // NOTE: updated to new serve mux instead of default one
	// mux := http.NewServeMux()

	// // sets up the main GraphQL endpoint where clients can send queries and mutations.
	// mux.Handle("/graphql", handler.New(s.ToExecutableSchema()))
	// // provides a web-based GraphQL Playground interface for easy testing and exploration of the GraphQL API.
	// mux.Handle("/playground", playground.Handler("play", "/graphql"))
	//
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.Handle("/graphql", handler.GraphQL(s.ToExecutableSchema()))
	http.Handle("/playground", handler.Playground("play", "/graphql"))

	// // Replace deprecated handler.GraphQL with handler.New
	// http.Handle("/graphql", handler.New(s.ToExecutableSchema()))
	// http.Handle("/playground", playground.Handler("play", "/graphql"))

	// Run server
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
