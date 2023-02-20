// Copyright (C) 2023  Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
