package main

// scraper is a long running job, scraper func is going to run in the background as our server run.
import (
	"time"

	"github.com/ankitmukhia/rssagg/internal/database"
)

func startScraping(database *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	//?
	/* 1. log message that says scrapping strts */
}
