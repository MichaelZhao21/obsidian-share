package main

import (
	"fmt"
	"net/http"
	"obsidianshare/src"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const port = 8080

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

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
	r.Route("/{id}", func(r chi.Router) {
	})

	// Admin routes
	r.Route("/admin", func(r chi.Router) {

	})

	// Start server
	fmt.Printf("Server is running on port %d\nGo to http://localhost:%d\n", port, port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
