package main

import (
	"log"
	"time"

	"github.com/haroonalbar/go-grpc-graphql-microservices/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	// Process environment variables and populate the config struct
	// Empty string means no prefix is used for env variables
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("error processing envconfig: %v", err)
	}
	var r catalog.Repository
	// Retry forever with 2-second intervals until successfully connecting to the database
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		// Attempt to create new Elasticsearch repository connection
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	// Ensure repository connection is closed when main function exits
	defer r.Close()
	log.Println("Listening on port: 8080")

	// Create new catalog service with the repository
	s := catalog.NewService(r)
	// Start gRPC server on port 8080
	log.Fatal(catalog.ListenGRPC(s, 8080))
}
