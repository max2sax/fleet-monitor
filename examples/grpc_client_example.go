package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/max2sax/fleet-monitor/grpc/client"
)

func main() {
	// Example of using the gRPC client

	// Connect to gRPC server
	c, err := client.NewFleetMonitorClient("localhost:6734")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer c.Close()

	ctx := context.Background()

	// Example 1: Send a heartbeat
	fmt.Println("=== Sending Heartbeat ===")
	heartbeatMsg, err := c.SendHeartbeat(ctx, "device-001", time.Now())
	if err != nil {
		log.Printf("SendHeartbeat error: %v", err)
	} else {
		fmt.Printf("Response: %s\n\n", heartbeatMsg)
	}

	// Example 2: Upload stats
	fmt.Println("=== Uploading Stats ===")
	uploadMsg, err := c.UploadStats(ctx, "device-001", 5000)
	if err != nil {
		log.Printf("UploadStats error: %v", err)
	} else {
		fmt.Printf("Response: %s\n\n", uploadMsg)
	}

	// Example 3: Get stats
	fmt.Println("=== Getting Stats ===")
	stats, err := c.GetStats(ctx, "device-001")
	if err != nil {
		log.Printf("GetStats error: %v", err)
	} else {
		fmt.Printf("Device ID: %s\n", stats.DeviceId)
		fmt.Printf("Avg Upload Time: %s\n", stats.AvgUploadTime)
		fmt.Printf("Uptime: %f\n", stats.Uptime)
	}

	// Example 4: Multiple devices
	fmt.Println("\n=== Multiple Devices ===")
	devices := []string{"device-002", "device-003", "device-004"}
	for _, deviceID := range devices {
		msg, err := c.SendHeartbeat(ctx, deviceID, time.Now())
		if err != nil {
			log.Printf("Error for %s: %v", deviceID, err)
		} else {
			fmt.Printf("%s: %s\n", deviceID, msg)
		}
	}
}
