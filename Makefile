.PHONY: proto
proto:
	@echo "Generating proto files..."
	protoc --go_out=. --go-grpc_out=. proto/fleet.proto

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o fleet-monitor

.PHONY: clean
clean:
	rm -rf grpc/pb/*.pb.go

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make proto   - Generate proto files"
	@echo "  make run     - Run the application"
	@echo "  make build   - Build the application"
	@echo "  make clean   - Clean generated files"
