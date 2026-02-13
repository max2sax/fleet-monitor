package main

import (
	"fmt"
	"net/http"

	"github.com/max2sax/fleet-monitor/api"
	"github.com/max2sax/fleet-monitor/storage"
)

func main() {
	// Initialize storage
	store := storage.NewStorage()

	// Initialize HTTP server
	server := &http.Server{
		Addr: ":6733",
	}

	// Initialize API
	deviceAPI := api.NewAPI(store, server).
		RegisterRoutes()

	// Start server
	fmt.Println("Server starting on :6733")
	if err := deviceAPI.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
