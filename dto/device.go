package dto

import "time"

type DeviceStatUpdate struct {
	DeviceId       string
	HeartbeatTime  *int64
	UploadDuration *int64
}

type DeviceStatDownload struct {
	DeviceId          string  `json:"device_id"`
	AverageUploadTime float64 `json:"avg_upload_time"`
	Uptime            float64 `json:"uptime"`
}

type DeviceHeartbeat struct {
	SentAt time.Time `json:"sent_at"`
}

type DeviceStatsUpload struct {
	SentAt     time.Time `json:"sent_at"`
	UploadTime int64     `json:"upload_time"`
}

type DeviceStatsDownload struct {
	AvgUploadTime string  `json:"avg_upload_time"` // returned as a time duration string. Eg: 5m10s
	Uptime        float64 `json:"uptime"`          // Uptime as a percentage. eg: 98.999
}
