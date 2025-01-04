package main

// scraper is a long running job, scraper func is going to run in the background as our server run.
import (
	"context"
	"database/sql"
	"log"
	"sync"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ankitmukhia/rssagg/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %v duration", concurrency, timeBetweenRequest)
	/* you can think of time.NewTicker() same as setInterval in js, not entirely */
	ticker := time.NewTicker(timeBetweenRequest)

	/* for range ticker.C, this will wait up-front, you try and run the build, the fire will wait for given ticker, then run/fire. */

	for ; ; <-ticker.C { // this will fire immediately
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error fetching feeds:", err)
			// this func should always be running, if we return here that will be prob.
			continue
		}
		// logic for, fetches each feed individually, at the same time.

		// we need synchronization mechanism
		// before understanding withgroup, we need to understand time.Sleep(), it can also use for same thing, wait for certain amount of time, until goroutines finishes their job. But time.Sleep() makes program wait for a fixed amount of time. whether or not the goroutine finishes its work. This can lead to errors if the taks takes longer or shorter then expected. For instance, if the tak is delayed, the progam might exit before it  finishes.

		// sync.WaitGroup() ensures the the program waits only until all goroutiens have completed, regardless of how long they take, It is synchronized with actual completion of tasks.
		// https://www.geeksforgeeks.org/using-waitgroup-in-golang/
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(), 
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,      
			Description: description,
			PublishedAt: time.Now(),
			Url: item.Link, 
			FeedID: feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue	
			}
			log.Println("Failed to create post:", err)
		}
	}
	log.Printf("feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
