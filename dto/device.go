package dto

import "time"

type DeviceStatUpdate struct {
	DeviceId       string
	HeartbeatTime  *int64
	UploadDuration *int64
}

type DeviceStatDownload struct {
	DeviceId          string  `json:"device_id"`
	AverageUploadTime string  `json:"avg_upload_time"`
	Uptime            float64 `json:"uptime"`
}

type DeviceHeartbeat struct {
	SentAt time.Time `json:"sent_at"`
}

type DeviceStatsUpload struct {
	UploadTime int64 `json:"upload_time"`
}
