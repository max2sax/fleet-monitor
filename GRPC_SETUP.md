# gRPC Setup for Fleet Monitor

This guide explains how to use the gRPC client and server that hooks into the existing HTTP API.

## Architecture

The fleet-monitor application now runs both HTTP and gRPC servers:
- **HTTP Server**: Running on `:6733` (existing REST API)
- **gRPC Server**: Running on `:6734` (new gRPC interface)

Both servers share the same underlying storage and business logic, allowing clients to interact with the API using either protocol.

## Supported gRPC Operations

The FleetMonitor gRPC service provides three main operations:

### 1. SendHeartbeat
Sends a device heartbeat to track device connectivity.

```go
SendHeartbeat(ctx, &pb.HeartbeatRequest{
    DeviceId: "device-123",
    SentAt: timestamppb.Now(),
})
```

### 2. UploadStats
Uploads device statistics including upload duration.

```go
UploadStats(ctx, &pb.UploadStatsRequest{
    DeviceId: "device-123",
    UploadTime: 5000, // milliseconds
})
```

### 3. GetStats
Retrieves statistics for a specific device.

```go
GetStats(ctx, &pb.GetStatsRequest{
    DeviceId: "device-123",
})
```

## Building Proto Files

If you need to regenerate the proto files:

```bash
make proto
```

This requires `protoc` to be installed. Install it with:
```bash
# macOS
brew install protobuf

# Linux
apt-get install protobuf-compiler

# Or from source
# https://github.com/protocolbuffers/protobuf
```

You'll also need the Go code generation plugins:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Running the Application

```bash
go run main.go
```

Both servers will start:
- HTTP: http://localhost:6733
- gRPC: localhost:6734

## Using the gRPC Client

### Connecting to the Server

```go
import (
    "github.com/max2sax/fleet-monitor/grpc/client"
)

// Create client
client, err := client.NewFleetMonitorClient("localhost:6734")
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Use client
ctx := context.Background()
response, err := client.SendHeartbeat(ctx, "device-123", time.Now())
if err != nil {
    log.Fatal(err)
}
```

## File Structure

```
fleet-monitor/
├── proto/
│   └── fleet.proto                 # Protocol Buffer definitions
├── grpc/
│   ├── pb/
│   │   ├── fleet.pb.go             # Generated message types
│   │   └── fleet_grpc.pb.go        # Generated gRPC service
│   ├── server.go                   # gRPC server implementation
│   └── client/
│       └── client.go               # gRPC client wrapper
├── main.go                         # Updated to run both servers
└── Makefile                        # Build and proto targets
```

## Protocol Buffer Definition

The gRPC service is defined in `proto/fleet.proto`:

```protobuf
service FleetMonitor {
  rpc SendHeartbeat(HeartbeatRequest) returns (HeartbeatResponse);
  rpc UploadStats(UploadStatsRequest) returns (UploadStatsResponse);
  rpc GetStats(GetStatsRequest) returns (GetStatsResponse);
}
```

## Implementation Details

### Server Implementation (grpc/server.go)

The gRPC server implements the FleetMonitorServer interface and hooks into the existing storage layer:

```go
type FleetMonitorServer struct {
    pb.UnimplementedFleetMonitorServer
    storage *storage.Storage
}

func (s *FleetMonitorServer) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
    // Converts gRPC request to internal DTOs
    // Calls storage.UpdateDeviceStats()
    // Returns gRPC response
}
```

### Client Wrapper (grpc/client/client.go)

The client wrapper provides a convenient interface to the gRPC service:

```go
type FleetMonitorClient struct {
    conn   *grpc.ClientConn
    client pb.FleetMonitorClient
}

func (c *FleetMonitorClient) SendHeartbeat(ctx context.Context, deviceID string, sentAt time.Time) (string, error)
func (c *FleetMonitorClient) UploadStats(ctx context.Context, deviceID string, uploadTime int64) (string, error)
func (c *FleetMonitorClient) GetStats(ctx context.Context, deviceID string) (*pb.GetStatsResponse, error)
```

## Testing the gRPC Server

Use `grpcurl` for testing:

```bash
# Install grpcurl
brew install grpcurl

# List available services
grpcurl -plaintext localhost:6734 list

# Call SendHeartbeat
grpcurl -plaintext -d '{"device_id": "test-device", "sent_at": "2026-02-21T00:00:00Z"}' localhost:6734 fleet.FleetMonitor/SendHeartbeat

# Call UploadStats
grpcurl -plaintext -d '{"device_id": "test-device", "upload_time": 1000}' localhost:6734 fleet.FleetMonitor/UploadStats

# Call GetStats
grpcurl -plaintext -d '{"device_id": "test-device"}' localhost:6734 fleet.FleetMonitor/GetStats
```

## Extending the Service

To add new gRPC methods:

1. Add the RPC definition to `proto/fleet.proto`
2. Implement the method in `grpc/server.go`
3. Add a wrapper method to `grpc/client/client.go`
4. Regenerate proto files: `make proto`

## Integration with Existing API

Both the HTTP and gRPC servers use the same underlying `storage.Storage` instance, so:
- Changes made through gRPC are immediately visible via HTTP
- Changes made through HTTP are immediately visible via gRPC
- Both protocols operate on the same data consistency model
