package main

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Welcome to your personal SQL query thing. We are launching immediately...")

	// Routine tells us when it's done
	done := make(chan int)

	// Total number of routines that finished
	total := 0

	// How many routines to run
	concurrency := 20

	// Create a pool of connections, "concurrency" in total so we have a connection per routine
	pool := pg.Connect(&pg.Options{
		User:     "postgres", // Change this to your user
		// Password: "", // Add your password here
		Addr:     "localhost:5432",
		PoolSize: concurrency, // Up to 20 goroutines
	})

	// Launch our little workers...
	for i := 0; i < concurrency; i++ {
		go sql(pool, done, i)
	}

	// Wait for them to be all done
	for {
		total += <-done
		if total >= concurrency {
			return
		}
	}
}

func sql(pool *pg.DB, done chan int, name int) {
	// Wait a random amount of time before starting so we randomize a little bit
	wait := time.Duration(rand.Intn(10000) * 60 * 1000)

	// Start with a wait
	time.Sleep(wait)

	// Just to know how far along we are
	var total int

	// Run 1024 queries
	for i := 0; i < 4096; i++ {
		// Change this to whatever you want
		query := "SELECT 1"

		// Return datatype, should be just one column like the results of:
		// COUNT(*) AS "count
		// MAX(id) AS "max"
		// etc...
		var n int

		// Checkout the connection from the pool and query the DB
		_, err := pool.QueryOne(pg.Scan(&n), query)

		// Any problems?
		if err != nil {
			fmt.Printf("[routine %d]: %s\n", name, err)
		} else {
			// We're good, let's count how far along we are
			total += n

			// Time to give the user an update ;)
			if total%40 == 0 {
				fmt.Printf("[routine %d]: ran %d queries, %d queries left\n", name, total, 4096-total)
			}
		}
		// Sleep for a bit
		time.Sleep(wait)
	}
	done <- 1
}
