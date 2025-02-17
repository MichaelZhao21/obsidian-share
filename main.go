package main

import (
	"fmt"
	"net/http"
	"obsidianshare/src"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load from dotenv for dev
	godotenv.Load()

	// Make sure all required environment variables are set
	src.CheckEnvs()

	// Create a new router
	r := chi.NewRouter()

	// Connect to MongoDB
	db := src.Connect()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(src.WithMongoDB(db))

	// Serve static files (including favicon.ico) from the "public" directory
	fs := http.FileServer(http.Dir("./public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	// ### Routes ###

	// Favicon route
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/favicon.ico")
	})

	// Heartbeat route
	r.Get("/", src.Index)

	// File routes
	r.Get("/{id}", src.GetFile)

	// Admin routes
	r.Route("/admin", func(r chi.Router) {

	})
	r.Post("/admin/pull", src.PullFiles)

	// Get the port from the environment variables
	port := os.Getenv("PORT")

	// Start server
	fmt.Printf("Server is running on port %s\nGo to http://localhost:%s\n", port, port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
