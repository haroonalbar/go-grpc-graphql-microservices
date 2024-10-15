package main

import (
	"log"
	"time"

	"github.com/haroonalbar/go-grpc-graphql-microservices/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

// NOTE:
// So in main we are populating cfg of type Config using envconfig
// To get the db url to connect to the postgresdb
// After we get a new service from account and start listening to grpc service on port 8080

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Error processing env: %v", err)
	}

	var r account.Repository
	// FIX: Later
	//  some randome package called retry with 8 stars says it's archived so must be depricated
	//  should look into it
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		// connect to db
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	// close db
	defer r.Close()

	log.Println("Listening on port 8080")
	// get service
	s := account.NewService(r)
	// start grpc server on 8080
	log.Fatal(account.ListenGRPC(s, 8080))
}
