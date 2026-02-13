package models

import "time"

type Device struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DeviceStats struct {
	DeviceID                   string `json:"deviceId"`
	AverageUploadTimeNS        int64  `json:"avg_upload_time"`
	NumberOfUploads            int64
	NumberOfHeartBeats         int64
	FirstHeartBeat             int64
	LastHeartBeat              int64
	CumulativeHeartBeatMinutes int64
}
