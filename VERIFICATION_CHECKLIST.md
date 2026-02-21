# gRPC Implementation Verification Checklist

Use this checklist to verify your gRPC integration is complete and working correctly.

## ✅ Files Created

- [x] `proto/fleet.proto` - Service definitions
- [x] `grpc/server.go` - gRPC server implementation
- [x] `grpc/client/client.go` - Client wrapper
- [x] `grpc/pb/fleet.pb.go` - Generated message types
- [x] `grpc/pb/fleet_grpc.pb.go` - Generated service stubs
- [x] `main.go` - Updated with gRPC server
- [x] `go.mod` - Added gRPC dependencies
- [x] `Makefile` - Build targets
- [x] `examples/grpc_client_example.go` - Example client
- [x] `GRPC_SETUP.md` - Detailed setup guide
- [x] `GRPC_QUICKSTART.md` - Quick start guide
- [x] `IMPLEMENTATION_SUMMARY.md` - Overview

## ✅ Dependencies

Verify gRPC and Protobuf are in go.mod:

```bash
grep "google.golang.org/grpc\|google.golang.org/protobuf" go.mod
```

Expected output:
```
google.golang.org/grpc v1.64.1
google.golang.org/protobuf v1.34.2
```

## ✅ Project Structure

Verify the new directories exist:

```bash
# These should all exist
ls -d proto/
ls -d grpc/
ls -d grpc/pb/
ls -d grpc/client/
ls -d examples/
```

## ✅ Build Checks

1. Download dependencies:
   ```bash
   go mod tidy
   ```

2. Verify the project compiles:
   ```bash
   go build -o fleet-monitor
   ```

3. Check for compilation errors:
   ```bash
   go build ./...
   ```

## ✅ Server Integration

Verify main.go contains:

```bash
# Check for gRPC imports
grep "grpc\|pb\|net.Listen" main.go | head -5

# Expected to see:
# - "net"
# - "net.Listen"
# - "google.golang.org/grpc"
# - pb imports
# - grpc.NewFleetMonitorServer
```

## ✅ Dual Server Launch

1. Start the application:
   ```bash
   go run main.go
   ```

2. Verify both servers start:
   - Look for: `HTTP Server starting on :6733`
   - Look for: `gRPC Server starting on :6734`

## ✅ Client Connection Test

Create a simple test file to verify connectivity:

```go
package main

import (
    "log"
    "github.com/max2sax/fleet-monitor/grpc/client"
)

func main() {
    c, err := client.NewFleetMonitorClient("localhost:6734")
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()
    log.Println("✓ Connected to gRPC server!")
}
```

Run it (in another terminal while server is running):
```bash
go run test_client.go
# Expected: "✓ Connected to gRPC server!"
```

## ✅ gRPC Method Tests

With server running, test each method:

### Test 1: SendHeartbeat
```bash
grpcurl -plaintext \
  -d '{"device_id":"test-device","sent_at":"2026-02-21T10:00:00Z"}' \
  localhost:6734 fleet.FleetMonitor/SendHeartbeat
```

Expected response:
```json
{
  "message": "heartbeat received"
}
```

### Test 2: UploadStats
```bash
grpcurl -plaintext \
  -d '{"device_id":"test-device","upload_time":1000}' \
  localhost:6734 fleet.FleetMonitor/UploadStats
```

Expected response:
```json
{
  "message": "stats uploaded successfully"
}
```

### Test 3: GetStats
```bash
grpcurl -plaintext \
  -d '{"device_id":"test-device"}' \
  localhost:6734 fleet.FleetMonitor/GetStats
```

Expected response:
```json
{
  "device_id": "test-device",
  "avg_upload_time": "1000",
  "uptime": 0
}
```

## ✅ Code Quality

1. Check for formatting issues:
   ```bash
   go fmt ./...
   ```

2. Run linter (if gofmt not enough):
   ```bash
   go vet ./...
   ```

3. Check module cleanliness:
   ```bash
   go mod verify
   ```

## ✅ Documentation

Verify documentation files are present:

```bash
ls -1 *.md | grep -i grpc
```

Expected files:
- GRPC_SETUP.md
- GRPC_QUICKSTART.md
- IMPLEMENTATION_SUMMARY.md

## ✅ HTTP API Still Works

Make sure existing HTTP API still works:

```bash
# With server running, test HTTP endpoint
curl -X GET http://localhost:6733/api/v1/devices/test-device/stats
```

Should work normally (may return not found if device doesn't exist)

## ✅ Data Consistency

Verify both protocols share data:

1. In Terminal 1 - Start server:
   ```bash
   go run main.go
   ```

2. In Terminal 2 - Use gRPC to send heartbeat:
   ```bash
   grpcurl -plaintext \
     -d '{"device_id":"sync-test","sent_at":"2026-02-21T10:00:00Z"}' \
     localhost:6734 fleet.FleetMonitor/SendHeartbeat
   ```

3. In Terminal 3 - Query via HTTP:
   ```bash
   curl http://localhost:6733/api/v1/devices/sync-test/stats
   # Should see the device with heartbeat data
   ```

## ✅ Performance Check

The gRPC server should not impact HTTP server:

1. Load test HTTP:
   ```bash
   ab -n 100 http://localhost:6733/api/v1/devices/test/stats
   ```

2. Load test gRPC:
   ```bash
   # Use ghz or similar gRPC load testing tool
   ghz --insecure -d '{"device_id":"test"}' \
       -n 100 localhost:6734 fleet.FleetMonitor/GetStats
   ```

Both should perform independently

## ✅ Cleanup (Optional)

Remove test files:
```bash
rm -f test_client.go
```

## 🎉 All Systems Go!

If all checks pass, your gRPC integration is complete and working correctly!

### Quick Reference - Common Commands

```bash
# Start server
go run main.go

# Run example client
go run examples/grpc_client_example.go

# Regenerate proto files (requires protoc)
make proto

# Build binary
make build

# Test HTTP API
curl http://localhost:6733/api/v1/devices/test/stats

# Test gRPC (requires grpcurl)
grpcurl -plaintext localhost:6734 list
```

## Troubleshooting

If something doesn't work:

1. Check server is running: `lsof -i :6734`
2. Check logs in console for error messages
3. Verify go.mod has correct versions
4. Try `go mod tidy && go build ./...`
5. See GRPC_SETUP.md for detailed troubleshooting
