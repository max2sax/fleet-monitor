package dto

type DeviceStatUpdate struct {
	DeviceId       string
	HeartbeatTime  *int64
	UploadDuration *int64
}

type DeviceStatDownload struct {
	DeviceId          string  `json: "device_id"`
	AverageUploadTime float64 `json: "avg_upload_time"`
	Uptime            float64 `json: "uptime"`
}
