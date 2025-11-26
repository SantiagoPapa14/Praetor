// main.go
package main

import (
	"Praetor/internal/db"
	"Praetor/internal/handlers"
	"Praetor/internal/middleware"
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

	// Repositories
	phraseRepo := repositories.NewPhraseRepository(database)
	sessionRepo := repositories.NewSessionRepository(database)

	// Handlers
	phrases := &handlers.PhraseHandler{Repository: phraseRepo}
	dashboard := &handlers.DashboardHandler{}
	authentication := &handlers.AuthenticationHandler{Repository: sessionRepo}

	// Router
	router := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("GET /", phrases.Page)
	router.HandleFunc("GET /authenticate", authentication.LoginPage)
	router.HandleFunc("GET /login", authentication.Login)
	router.HandleFunc("GET /phrases", phrases.List)
	router.HandleFunc("POST /phrases", phrases.Add)
	router.HandleFunc("DELETE /phrases/{id}", phrases.Delete)

	protected := middleware.AuthMiddleware(sessionRepo, http.HandlerFunc(dashboard.Page))
	router.Handle("GET /dashboard", protected)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	log.Println("🚀 Server running at http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
