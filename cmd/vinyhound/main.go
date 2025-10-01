package main

import (
	"log"
	"net/http"
	"os"

	"vinyhound/internal/app"
)

func main() {
	store := app.NewStore()

	if err := store.CreateUser("demo", "demo123", []string{
		"Welcome to Vinyhound!",
		"Start by customizing your personal playlist.",
	}); err != nil {
		log.Fatalf("bootstrap demo user: %v", err)
	}

	server := app.NewServer(store)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("API available at http://localhost%v", addr)
	if err := http.ListenAndServe(addr, server.Routes()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
