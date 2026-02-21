package grpc

import (
	"context"
	"errors"

	"github.com/max2sax/fleet-monitor/dto"
	pb "github.com/max2sax/fleet-monitor/grpc/pb"
	"github.com/max2sax/fleet-monitor/storage"
)

type FleetMonitorServer struct {
	pb.UnimplementedFleetMonitorServer
	storage *storage.Storage
}

func NewFleetMonitorServer(store *storage.Storage) *FleetMonitorServer {
	return &FleetMonitorServer{
		storage: store,
	}
}

func (s *FleetMonitorServer) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	if req.DeviceId == "" {
		return nil, errors.New("device_id is required")
	}
	if req.SentAt == nil {
		return nil, errors.New("sent_at is required")
	}

	heartbeatTime := req.SentAt.AsTime().Unix()
	update := dto.DeviceStatUpdate{
		DeviceId:      req.DeviceId,
		HeartbeatTime: &heartbeatTime,
	}

	err := s.storage.UpdateDeviceStats(&update)
	if errors.Is(err, &dto.ErrorNotFound{}) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &pb.HeartbeatResponse{
		Message: "heartbeat received",
	}, nil
}

func (s *FleetMonitorServer) UploadStats(ctx context.Context, req *pb.UploadStatsRequest) (*pb.UploadStatsResponse, error) {
	if req.DeviceId == "" {
		return nil, errors.New("device_id is required")
	}

	update := dto.DeviceStatUpdate{
		DeviceId:       req.DeviceId,
		UploadDuration: &req.UploadTime,
	}

	err := s.storage.UpdateDeviceStats(&update)
	if errors.Is(err, &dto.ErrorNotFound{}) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &pb.UploadStatsResponse{
		Message: "stats uploaded successfully",
	}, nil
}

func (s *FleetMonitorServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	if req.DeviceId == "" {
		return nil, errors.New("device_id is required")
	}

	dev, err := s.storage.GetDeviceStats(req.DeviceId)
	if errors.Is(err, &dto.ErrorNotFound{}) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &pb.GetStatsResponse{
		DeviceId:      dev.DeviceId,
		AvgUploadTime: dev.AverageUploadTime,
		Uptime:        dev.Uptime,
	}, nil
}
