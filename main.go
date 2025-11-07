package main

import (
	"log"
	"net/http"

	"Praetor/internal/db"
	"Praetor/internal/handlers"
)

func main() {
	database, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	router := http.NewServeMux()

	phrases := &handlers.PhraseHandler{DB: database}

	router.HandleFunc("GET /", phrases.List)
	router.HandleFunc("POST /phrases", phrases.Add)
	router.HandleFunc("DELETE /phrases/{id}", phrases.Delete)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	log.Println("🚀 Server running at http://localhost:8080")
	server.ListenAndServe()
}
