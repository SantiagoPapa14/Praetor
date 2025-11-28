package main

import (
	"Praetor/internal/app"
	"Praetor/internal/db"
	"Praetor/internal/handlers"
	"Praetor/internal/middleware"
	"Praetor/internal/repositories"
	"context"
	"log"
	"net/http"

	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()

	database, err := db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// App
	application := &app.App{}
	application.Repos.Session = repositories.NewSessionRepository(database)
	application.Repos.User = repositories.NewUserRepository(database)
	application.Repos.Docker = repositories.NewDockerRepository(cli, ctx)

	// Handlers
	dashboard := &handlers.DashboardHandler{}
	authentication := &handlers.AuthenticationHandler{App: application}

	// Router
	router := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("GET /authenticate", authentication.LoginPage)
	router.HandleFunc("POST /auth/login", authentication.Login)
	router.HandleFunc("POST /auth/register", authentication.Register)
	router.HandleFunc("POST /auth/logout", authentication.Logout)

	// Temp
	tempDockerHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		application.Repos.Docker.GetContainers()
	})
	router.HandleFunc("GET /containers", tempDockerHandlerFunc)

	protectedDashboard := middleware.AuthMiddleware(application, http.HandlerFunc(dashboard.Page))
	router.Handle("GET /", protectedDashboard)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	log.Println("🚀 Server running at http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
