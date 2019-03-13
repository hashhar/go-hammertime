package main

import (
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hashhar/go-hammertime/pkg"
	"github.com/hashhar/go-hammertime/pkg/database"
	"github.com/hashhar/go-hammertime/pkg/models"
)

// Command line flags
var (
	concurrency     = flag.Int("concurrency", 10, "Number of worker goroutines to use")
	numItems        = flag.Int("items", 10000, "Number of items to create")
	batchSize       = flag.Int("batch-size", 1000, "Number of items to create in single query")
	progressSeconds = flag.Int("progress", 60, "Interval in seconds to report progress")
	dbHost          = flag.String("host", "localhost", "Hostname of the Postgres server")
	dbPort          = flag.Int("port", 5432, "Port of the Postgres server")
	dbName          = flag.String("database", "godamqadb_benchmark", "Database to use for benchmarking")
)

var (
	logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
)

func main() {
	flag.Parse()
	startTime := time.Now()

	defer func() {
		currentItems := database.Count()
		logger.Printf(
			"Processed %d items of %d (%f qps)\n",
			currentItems,
			*numItems,
			2*float64(currentItems)/(time.Since(startTime).Seconds()), // Multiplied by 2 since it's an insert + update
		)
		database.Close()
	}()

	database.InitDB(*dbHost, *dbPort, *dbName)

	// Progress reporting
	ticker := time.NewTicker(time.Duration(*progressSeconds) * time.Second)
	go func() {
		for tick := range ticker.C {
			currentItems := database.Count()
			logger.Printf(
				"Processed %d items of %d (%f qps)\n",
				currentItems,
				*numItems,
				2*float64(currentItems)/(tick.Sub(startTime).Seconds()), // Multiplied by 2 since it's an insert + update
			)
		}
	}()

	// Channel to send work on. Capacity can be kept smaller than numItems since
	// is is unreasonable to expect out worker to keep up with our work.
	work := make(chan models.Work, *numItems)
	go pkg.SendWork(work, *numItems)

	wg := sync.WaitGroup{}
	// Start concurrency number of goroutines
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go hammerTime(work, &wg)
	}
	wg.Wait()
}

func hammerTime(work chan models.Work, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()

	if *batchSize != 1 {
		batch := make([]models.Work, 0, *batchSize)
		for w := range work {
			if len(batch) < *batchSize {
				batch = append(batch, w)
			} else {
				database.BatchCreateAndUpdateStatus(batch)
				batch = batch[:0]
				batch = append(batch, w)
			}
		}
		database.BatchCreateAndUpdateStatus(batch)
	} else {
		for w := range work {
			database.CreateAndUpdateStatus(w)
		}
	}
}
