package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/max2sax/fleet-monitor/api"
	"github.com/max2sax/fleet-monitor/grpc"
	pb "github.com/max2sax/fleet-monitor/grpc/pb"
	"github.com/max2sax/fleet-monitor/storage"
	grpclib "google.golang.org/grpc"
)

func main() {
	// Initialize storage
	store := storage.NewStorage()

	// Initialize HTTP server
	httpServer := &http.Server{
		Addr: ":6733",
	}

	// Initialize API
	deviceAPI := api.NewAPI(store, httpServer).
		RegisterRoutes()

	// Start HTTP server in a goroutine
	go func() {
		fmt.Println("HTTP Server starting on :6733")
		if err := deviceAPI.Start(); err != nil {
			fmt.Printf("HTTP Server error: %v\n", err)
		}
	}()

	// Initialize and start gRPC server
	grpcListener, err := net.Listen("tcp", ":6734")
	if err != nil {
		fmt.Printf("Failed to listen for gRPC: %v\n", err)
		return
	}

	grpcServer := grpclib.NewServer()
	fleetMonitorServer := grpc.NewFleetMonitorServer(store)
	pb.RegisterFleetMonitorServer(grpcServer, fleetMonitorServer)

	fmt.Println("gRPC Server starting on :6734")
	if err := grpcServer.Serve(grpcListener); err != nil {
		fmt.Printf("gRPC Server error: %v\n", err)
	}
}
