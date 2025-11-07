// main.go
package main

import (
	"Praetor/internal/db"
	"Praetor/internal/handlers"
	"Praetor/internal/repositories"
	"log"
	"net/http"
)

func main() {
	database, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Create repository
	phraseRepo := repositories.NewPhraseRepository(database)

	// Create handler with repository
	phrases := &handlers.PhraseHandler{Repository: phraseRepo}

	router := http.NewServeMux()
	router.HandleFunc("GET /", phrases.Page)
	router.HandleFunc("GET /phrases", phrases.List)
	router.HandleFunc("POST /phrases", phrases.Add)
	router.HandleFunc("DELETE /phrases/{id}", phrases.Delete)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	log.Println("🚀 Server running at http://localhost:8080")
	server.ListenAndServe()
}
