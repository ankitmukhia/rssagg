package main

import (
	"log"
	"os"
	"time"
	"database/sql"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	/* https://stackoverflow.com/questions/52789531/how-do-i-solve-panic-sql-unknown-driver-postgres-forgotten-import */
 	 _ "github.com/lib/pq"

	"github.com/ankitmukhia/rssagg/internal/database"
)

//NOTE:Opened connection to db, and stored it in the state struct. Now I can perform all regular db queries.
type state struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found")
	}

	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Fatal("URL not found")
	}

	dbCon, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Unable to connect to db", err)
	}
	db := database.New(dbCon)
	dbQueries := state{
		DB: db,
	}

	go startScraping(db, 10, time.Minute) 

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	v1Router.Post("/users", dbQueries.handlerCreateUser)
	v1Router.Get("/users", dbQueries.middlewareAuth(dbQueries.handlerGetUser))

	v1Router.Post("/feeds", dbQueries.middlewareAuth(dbQueries.handlerCreateFeed))
	v1Router.Get("/feeds", dbQueries.handlerGetFeeds)

	v1Router.Post("/follow_feed", dbQueries.middlewareAuth(dbQueries.handlerFollowFeed))
	v1Router.Get("/follow_feed", dbQueries.middlewareAuth(dbQueries.handlerGetFollowFeed))
	v1Router.Delete("/follow_feed/{feedFollowId}", dbQueries.middlewareAuth(dbQueries.handlerDeleteFollowFeed))
	v1Router.Get("/posts", dbQueries.middlewareAuth(dbQueries.handlerGetPostsForUser))

	/* attaches another http.handler as a subrouter: v1/healthz */
	r.Mount("/v1", v1Router)
	
	srv := &http.Server{
		Addr: ":" + port,
		Handler: r,
	}
	
	/* Code actually stops right here, and starts handaling http requests. Until we don't have any err, which will be handled by Fatal. */
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
