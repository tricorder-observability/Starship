package main

import (
	"flag"
	"log"
	"time"

	"github.com/tricorder/src/utils/http"
	"github.com/tricorder/src/utils/pg"
)

func main() {
	pgURL := flag.String("pg_url", "", "The URL of the Postgres server")
	genFrequency := flag.Int("gen_frequency", 1, "The frequency for per second")

	flag.Parse()

	client := pg.NewClient(*pgURL)
	if err := client.Connect(); err != nil {
		log.Fatalf("Could not connect to Postgres server at %s", *pgURL)
	}
	defer client.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 1 * *genFrequency times per second.
		for i := 0; i < *genFrequency; i++ {
			req := http.Gen()
			_ = client.WriteHTTPRequest(&req)
		}
	}
}
