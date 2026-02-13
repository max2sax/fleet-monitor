package storage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/max2sax/fleet-monitor/dto"
	"github.com/max2sax/fleet-monitor/models"
)

type messageWriteRequest struct {
	message *dto.DeviceStatUpdate
	result  chan error
}

type Storage struct {
	stats           map[string]models.DeviceStats // map[string][]models.Message
	deviceWriteChan chan messageWriteRequest
}

func NewStorage() *Storage {
	s := &Storage{
		deviceWriteChan: make(chan messageWriteRequest),
	}
	go s.deviceWriter()
	// init devices
	s.stats = make(map[string]models.DeviceStats)

	// read csv file with devices
	file, err := os.Open("devices.csv")
	if err != nil {
		log.Fatal(err)
	}
	// 2. Ensure the file is closed when the function returns
	defer file.Close()

	// 3. Create a scanner for the file
	scanner := bufio.NewScanner(file)

	// 4. Iterate through the file line by line
	for scanner.Scan() {
		id := scanner.Text() // Get the current line as a string
		fmt.Println("loading device: " + id)
		s.stats[id] = models.DeviceStats{
			DeviceID: id,
		}
	}

	// 5. Check for errors encountered during scanning
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return s
}

func (s *Storage) deviceWriter() {
	for req := range s.deviceWriteChan {
		updateRequest := req.message

		ds, ok := s.stats[updateRequest.DeviceId]
		if !ok {
			err := dto.ErrorNotFound{What: "device(" + updateRequest.DeviceId + ")"}
			req.result <- &err
			continue
		}
		if updateRequest.HeartbeatTime != nil {
			heartbeatSeconds := *updateRequest.HeartbeatTime
			ds.NumberOfHeartBeats++
			if ds.FirstHeartBeat == 0 {
				ds.FirstHeartBeat = heartbeatSeconds
			}
			ds.CumulativeHeartBeatMinutes = ((heartbeatSeconds) - ds.FirstHeartBeat) / 60
		}
		if updateRequest.UploadDuration != nil {
			// calculate new average
			ds.NumberOfUploads++
			diffCA := *updateRequest.UploadDuration - ds.AverageUploadTimeNS
			ds.AverageUploadTimeNS = ds.AverageUploadTimeNS + (diffCA / ds.NumberOfUploads)
		}

		s.stats[updateRequest.DeviceId] = ds
		req.result <- nil
	}
}

func (s *Storage) UpdateDeviceStats(message *dto.DeviceStatUpdate) error {
	result := make(chan error, 1)
	s.deviceWriteChan <- messageWriteRequest{
		message: message,
		result:  result,
	}
	return <-result
}

func (s *Storage) GetDeviceStats(deviceId string) (*dto.DeviceStatDownload, error) {
	ds, ok := s.stats[deviceId]
	if !ok {
		return nil, fmt.Errorf("%w", &dto.ErrorNotFound{What: "device(" + deviceId + ")"})
	}
	dur := time.Duration(ds.AverageUploadTimeNS)
	uptime := (float64(ds.NumberOfHeartBeats) / float64(ds.CumulativeHeartBeatMinutes)) * 100.0
	dsd := dto.DeviceStatDownload{
		DeviceId:          deviceId,
		Uptime:            uptime,
		AverageUploadTime: dur.String(),
	}
	return &dsd, nil
}
