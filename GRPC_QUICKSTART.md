# Quick Start: gRPC Client Setup

## Prerequisites

Ensure you have Go 1.25+ installed and the required modules:

```bash
go mod download
```

## Starting the Server

The application now runs both HTTP and gRPC servers:

```bash
go run main.go
```

Both servers start automatically:
- **HTTP API**: http://localhost:6733
- **gRPC Service**: localhost:6734

## Using the gRPC Client

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/max2sax/fleet-monitor/grpc/client"
)

func main() {
    // Create connection
    c, err := client.NewFleetMonitorClient("localhost:6734")
    if err != nil {
        panic(err)
    }
    defer c.Close()

    // Send heartbeat
    ctx := context.Background()
    response, err := c.SendHeartbeat(ctx, "device-123", time.Now())
    if err != nil {
        panic(err)
    }
    fmt.Println(response) // "heartbeat received"
}
```

### Available Client Methods

```go
// Send device heartbeat
SendHeartbeat(ctx, deviceID string, sentAt time.Time) (string, error)

// Upload device stats
UploadStats(ctx, deviceID string, uploadTime int64) (string, error)

// Get device statistics
GetStats(ctx, deviceID string) (*pb.GetStatsResponse, error)
```

## Testing with grpcurl

If you have `grpcurl` installed:

```bash
# Send heartbeat
grpcurl -plaintext \
  -d '{"device_id":"test","sent_at":"2026-02-21T00:00:00Z"}' \
  localhost:6734 fleet.FleetMonitor/SendHeartbeat

# Upload stats
grpcurl -plaintext \
  -d '{"device_id":"test","upload_time":1000}' \
  localhost:6734 fleet.FleetMonitor/UploadStats

# Get stats
grpcurl -plaintext \
  -d '{"device_id":"test"}' \
  localhost:6734 fleet.FleetMonitor/GetStats
```

## Running the Example Client

A complete example is provided in `examples/grpc_client_example.go`:

```bash
# Update the import path and run
go run examples/grpc_client_example.go
```

## Project Structure

- **proto/fleet.proto** - gRPC service definitions
- **grpc/server.go** - gRPC server implementation
- **grpc/client/client.go** - Client wrapper for easy use
- **grpc/pb/** - Generated protobuf code
- **examples/grpc_client_example.go** - Usage example

## Data Flow

```
gRPC Client
    ↓
gRPC Server (grpc/server.go)
    ↓
Storage Layer (storage/storage.go)
    ↓
Device Data

HTTP Client
    ↓
HTTP API (api/api.go)
    ↓
Storage Layer (storage/storage.go)
    ↓
Device Data
```

Both protocols share the same storage, ensuring data consistency.

## Next Steps

See [GRPC_SETUP.md](GRPC_SETUP.md) for:
- Regenerating proto files
- Extending with new RPC methods
- Advanced configuration options
- Error handling
