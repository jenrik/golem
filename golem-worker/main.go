package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jenrik/golem"
	_ "github.com/jenrik/golem/action"
	"os"
	"strconv"
)

func main() {
	// Setup logging
	if env, ok := os.LookupEnv("ENVIRONMENT"); ok && env == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// Connection details
	storage, ok := os.LookupEnv("GOLEM_STORAGE")
	if !ok {
		fmt.Printf("Missing storage address. Set environment variable GOLEM_STORAGE\n")
		return
	}
	queue, ok := os.LookupEnv("GOLEM_QUEUE")
	if !ok {
		fmt.Printf("Missing queue address. Set environment variable GOLEM_QUEUE\n")
		return
	}
	err := golem.Connect(queue, storage)
	if err != nil {
		panic(err)
	}

	// Start workers
	var workers int
	sworkers, ok := os.LookupEnv("GOLEM_WORKERS")
	if !ok {
		workers = 1
	} else {
		iworkers, err := strconv.ParseInt(sworkers, 10, 0)
		if err != nil {
			panic(err)
		}
		workers = int(iworkers)
	}
	wg := golem.StartWorkers(workers)
	wg.Wait()
}
