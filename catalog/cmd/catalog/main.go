package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sriraksharao/go-microservice/catalog"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	log.Println("listening on port 8080")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))
}
