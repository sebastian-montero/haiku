package main

import (
	"log"
	"net/http"
	"os"

	"haiku/internal/app"
)

func main() {
	dataDir := os.Getenv("HAIKU_DATA")
	if dataDir == "" {
		dataDir = "./data"
	}
	addr := os.Getenv("HAIKU_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	store, err := app.NewStore(dataDir)
	if err != nil {
		log.Fatalf("init store: %v", err)
	}
	hub := app.NewHub(store)
	srv := app.NewServer(store, hub)

	log.Printf("haiku listening on %s (data dir: %s)", addr, dataDir)
	log.Fatal(http.ListenAndServe(addr, srv.Routes()))
}
