# Fleet Monitor

A lightweight Go API server for monitoring device heartbeats and collecting device statistics across a fleet of devices.

## Overview

Fleet Monitor is a REST API that tracks device availability and performance metrics. It manages device heartbeats and collects upload time statistics from devices in your fleet. The server stores device information in memory and logs events for monitoring.

## Features

- **Heartbeat Monitoring**: Track device availability through periodic heartbeat signals
- **Statistics Collection**: Collect and aggregate device upload time metrics
- **Device Management**: Load and manage a fleet of devices from a CSV file
- **Error Handling**: Comprehensive error responses with meaningful status codes
- **Request Logging**: Built-in middleware to log all incoming HTTP requests

## Getting Started

### Prerequisites

- Go 1.25.4 or higher
- A `devices.csv` file containing device IDs (one per line)

### Installation

```bash
# Clone the repository (or navigate to the project directory)
cd /Users/max2sax-proto/development/fleet-monitor

# Install dependencies
go mod download
```

### Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:6733`

Expected output:

```text
Server starting on :6733
loading device: <device_id>
...
```

## API Endpoints

All endpoints return JSON responses and use the following base URL: `http://localhost:6733`

### 1. Send Device Heartbeat

**Endpoint:** `POST /api/v1/devices/{device_id}/heartbeat`

Send a heartbeat signal to indicate a device is active.

**Parameters:**

- `device_id` (path parameter, string): Unique identifier of the device

**Request Body:**

```json
{
  "sent_at": "2024-02-13T10:30:00Z"
}
```

**Response:**

- **204 No Content**: Heartbeat successfully recorded
- **400 Bad Request**: Missing device_id or invalid JSON
- **404 Not Found**: Device ID not found in the system
- **500 Internal Server Error**: Server-side processing error

**Example:**

```bash
curl -X POST http://localhost:6733/api/v1/devices/device-001/heartbeat \
  -H "Content-Type: application/json" \
  -d '{"sent_at": "2024-02-13T10:30:00Z"}'
```

---

### 2. Upload Device Statistics

**Endpoint:** `POST /api/v1/devices/{device_id}/stats`

Submit upload time statistics from a device.

**Parameters:**

- `device_id` (path parameter, string): Unique identifier of the device

**Request Body:**

```json
{
  "upload_time": 1500000000
}
```

**Fields:**

- `upload_time` (int64): Upload duration in nanoseconds

**Response:**

- **204 No Content**: Statistics successfully recorded
- **400 Bad Request**: Missing device_id or invalid JSON
- **404 Not Found**: Device ID not found in the system
- **500 Internal Server Error**: Server-side processing error

**Example:**

```bash
curl -X POST http://localhost:6733/api/v1/devices/device-001/stats \
  -H "Content-Type: application/json" \
  -d '{"upload_time": 1500000000}'
```

---

### 3. Get Device Statistics

**Endpoint:** `GET /api/v1/devices/{device_id}/stats`

Retrieve aggregated statistics for a specific device.

**Parameters:**

- `device_id` (path parameter, string): Unique identifier of the device

**Response:**

- **200 OK**: Returns device statistics
- **400 Bad Request**: Missing device_id parameter
- **404 Not Found**: Device ID not found in the system
- **500 Internal Server Error**: Server-side processing error

**Response Body:**

```json
{
  "deviceId": "device-001",
  "avg_upload_time": "1500000000",
  "uptime": 98.5
}
```

**Example:**

```bash
curl -X GET http://localhost:6733/api/v1/devices/device-001/stats
```

---

## Data Models

### DeviceHeartbeat (Request)

```go
{
  "sent_at": "2024-02-13T10:30:00Z"  // ISO 8601 timestamp
}
```

### DeviceStatsUpload (Request)

```go
{
  "upload_time": 1500000000  // Duration in nanoseconds
}
```

### DeviceStatDownload (Response)

```go
{
  "device_id": "device-001",
  "avg_upload_time": "1500000000",
  "uptime": 98.5
}
```

## Testing

### Manual Testing with cURL

#### Test Heartbeat Endpoint

```bash
curl -X POST http://localhost:6733/api/v1/devices/device-001/heartbeat \
  -H "Content-Type: application/json" \
  -d '{"sent_at": "2024-02-13T10:30:00Z"}'
```

#### Test Stats Upload Endpoint

```bash
curl -X POST http://localhost:6733/api/v1/devices/device-001/stats \
  -H "Content-Type: application/json" \
  -d '{"upload_time": 1500000000}'
```

#### Test Stats Retrieval Endpoint

```bash
curl -X GET http://localhost:6733/api/v1/devices/device-001/stats
```

#### Test with Invalid Device ID

```bash
curl -X GET http://localhost:6733/api/v1/devices/invalid-device/stats
# Expected: 404 error
```

### Using Device Simulator

A device simulator binary is available for testing:

```bash
./device-simulator-mac-arm64
```

This simulator will send heartbeats and statistics from devices defined in `devices.csv`.

### Response Examples

**Successful Stats Retrieval:**

```text
HTTP/1.1 200 OK
Content-Type: application/json

{"deviceId":"device-001","avg_upload_time":"1500000000","uptime":98.5}
```

**Device Not Found:**

```text
HTTP/1.1 404 Not Found

device not found
```

**Bad JSON Request:**

```text
HTTP/1.1 400 Bad Request
Content-Type: application/json

{"msg":"Bad JSON: unexpected end of JSON input"}
```

## Configuration

### Device Inventory

The server loads device IDs from a `devices.csv` file in the project root directory. Each line should contain one device ID:

```text
device-001
device-002
device-003
device-004
```

### Server Port

The server runs on port `6733` by default. To change this, modify the `Addr` field in [main.go](main.go):

```go
server := &http.Server{
    Addr: ":6733",  // Change port here
}
```

## Project Structure

```text
fleet-monitor/
├── main.go                 # Server entry point
├── api/
│   ├── api.go             # HTTP handlers and routing
│   └── log.go             # Request logging middleware
├── dto/
│   ├── device.go          # Data transfer objects
│   └── errors.go          # Custom error types
├── models/
│   └── models.go          # Domain models
├── storage/
│   └── storage.go         # In-memory storage and persistence
├── devices.csv            # Device inventory
└── README.md              # This file
```

## Error Handling

The API provides descriptive error messages in JSON format:

```json
{
  "msg": "error description"
}
```

Common HTTP status codes:

- `200 OK`: Request successful
- `204 No Content`: Request successful, no response body
- `400 Bad Request`: Invalid request parameters or malformed JSON
- `404 Not Found`: Device not found in inventory
- `500 Internal Server Error`: Unexpected server error

## Logging

All HTTP requests are logged automatically by the request logging middleware in [api/log.go](api/log.go). Check the console output for request details.

### Building

```bash
go build -o fleet-monitor
```

## License

See the [LICENSE](LICENSE) file for details.
