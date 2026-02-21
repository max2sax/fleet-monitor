package client

import (
	"context"
	"time"

	pb "github.com/max2sax/fleet-monitor/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FleetMonitorClient struct {
	conn   *grpc.ClientConn
	client pb.FleetMonitorClient
}

func NewFleetMonitorClient(addr string) (*FleetMonitorClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &FleetMonitorClient{
		conn:   conn,
		client: pb.NewFleetMonitorClient(conn),
	}, nil
}

func (c *FleetMonitorClient) SendHeartbeat(ctx context.Context, deviceID string, sentAt time.Time) (string, error) {
	req := &pb.HeartbeatRequest{
		DeviceId: deviceID,
		SentAt:   timestamppb.New(sentAt),
	}

	resp, err := c.client.SendHeartbeat(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Message, nil
}

func (c *FleetMonitorClient) UploadStats(ctx context.Context, deviceID string, uploadTime int64) (string, error) {
	req := &pb.UploadStatsRequest{
		DeviceId:   deviceID,
		UploadTime: uploadTime,
	}

	resp, err := c.client.UploadStats(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Message, nil
}

func (c *FleetMonitorClient) GetStats(ctx context.Context, deviceID string) (*pb.GetStatsResponse, error) {
	req := &pb.GetStatsRequest{
		DeviceId: deviceID,
	}

	return c.client.GetStats(ctx, req)
}

func (c *FleetMonitorClient) Close() error {
	return c.conn.Close()
}
