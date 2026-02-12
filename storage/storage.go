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
	devices         map[string]models.Device      // map[string]*models.Room
	stats           map[string]models.DeviceStats // map[string][]models.Message
	deviceWriteChan chan messageWriteRequest
}

func NewStorage() *Storage {
	s := &Storage{
		deviceWriteChan: make(chan messageWriteRequest),
	}
	go s.deviceWriter()
	// init devices
	s.devices = make(map[string]models.Device)
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
		nowt := time.Now().UTC()
		s.devices[id] = models.Device{
			ID:        id,
			CreatedAt: nowt,
			UpdatedAt: nowt,
		}
	}

	// 5. Check for errors encountered during scanning
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.Fatal("testing error")
	return s
}

func (s *Storage) deviceWriter() {
	for req := range s.deviceWriteChan {
		updateRequest := req.message
		d, ok := s.devices[updateRequest.DeviceId]
		if !ok {
			req.result <- fmt.Errorf("device not found")
			continue
		}
		ds, ok := s.stats[updateRequest.DeviceId]
		if !ok {
			req.result <- fmt.Errorf("stats not found")
			continue
		}
		//TODO: Update device stats with incoming data
		nowTime := time.Now().UTC()
		d.UpdatedAt = nowTime
		if updateRequest.HeartbeatTime != nil {
			ds.NumberOfHeartBeats++
			if ds.LastHeartBeat == 0 {
				ds.LastHeartBeat = nowTime.Unix() * 60
			}
			ds.CumulativeHeartBeatTime = (nowTime.Unix() * 60) - ds.LastHeartBeat
			ds.LastHeartBeat = nowTime.Unix() * 60
			//calculate the diff
		}
		if updateRequest.UploadDuration != nil {
			// calculate new average
			ds.NumberOfUploads++
			diffCA := *updateRequest.UploadDuration - ds.AverageUploadTime
			ds.AverageUploadTime = ds.AverageUploadTime + (diffCA / ds.NumberOfUploads)
		}

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

func (s *Storage) GetDevice(deviceId string) (*models.Device, error) {
	room, ok := s.devices[deviceId]
	if !ok {
		return nil, fmt.Errorf("device not found")
	}
	return &room, nil
}

func (s *Storage) GetAllDevices() []models.Device {
	var rooms []models.Device
	for _, v := range s.devices {
		rooms = append(rooms, v)
	}
	return rooms
}
