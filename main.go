package main

import (
	"log"
	"os"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found")
	}

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
