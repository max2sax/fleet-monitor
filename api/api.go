package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/max2sax/fleet-monitor/dto"
	"github.com/max2sax/fleet-monitor/storage"
)

type ErrorResponse struct {
	Msg string `json:"msg"`
}

type API struct {
	storage *storage.Storage
	mux     *http.ServeMux
	server  *http.Server
}

func NewAPI(store *storage.Storage, server *http.Server) *API {
	mux := http.NewServeMux()
	server.Handler = mux
	return &API{
		storage: store,
		mux:     mux,
		server:  server,
	}
}

func (a *API) RegisterRoutes() *API {
	a.mux.HandleFunc("POST /devices/{device_id}/heartbeat", a.heartbeatHandler)
	a.mux.HandleFunc("POST /devices/{device_id}/stats", a.uploadDeviceStats)
	a.mux.HandleFunc("GET /devices/{device_id}/stats", a.getDeviceStats)
	return a
}

func (a *API) Start() error {
	return a.server.ListenAndServe()
}

func (a *API) heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	deviceId := r.PathValue("device_id")
	if deviceId == "" {
		http.Error(w, "Device ID required", http.StatusBadRequest)
		return
	}
	var req dto.DeviceHeartbeat
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	timeNow := time.Now().UTC().Unix()
	update := dto.DeviceStatUpdate{
		DeviceId:      deviceId,
		HeartbeatTime: &timeNow,
	}
	// if device is not found then return a 404 with ErrorResponse and msg missing
	// if there is some other error return a 500 with ErrorResponse indicating error
	err := a.storage.UpdateDeviceStats(&update)
	if errors.Is(err, &dto.ErrorNotFound{}) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "unable to load device stats", http.StatusInternalServerError)
		return
	}
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(room)
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) uploadDeviceStats(w http.ResponseWriter, r *http.Request) {
	deviceId := r.PathValue("device_id")
	if deviceId == "" {
		http.Error(w, "Device ID required", http.StatusBadRequest)
		return
	}
	var req dto.DeviceStatsUpload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	update := dto.DeviceStatUpdate{
		DeviceId:       deviceId,
		UploadDuration: &req.UploadTime,
	}
	// if device is not found then return a 404 with ErrorResponse and msg missing
	// if there is some other error return a 500 with ErrorResponse indicating error
	err := a.storage.UpdateDeviceStats(&update)
	if errors.Is(err, &dto.ErrorNotFound{}) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "unable to load device stats", http.StatusInternalServerError)
		return
	}
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(room)
	w.WriteHeader(http.StatusNoContent)

	// TODO: call storage layer to update device stats
	// if device is not found then return a 404 with ErrorResponse and msg missing
	// if there is some other error return a 500 with ErrorResponse indicating error

	w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(msg)
}

func (a *API) getDeviceStats(w http.ResponseWriter, r *http.Request) {
	deviceId := r.PathValue("device_id")
	if deviceId == "" {
		http.Error(w, "device id required", http.StatusBadRequest)
		return
	}

	// TODO:
	// if device is not found then return a 404 with ErrorResponse and msg missing
	// if there is some other error return a 500 with ErrorResponse indicating error
	dev, err := a.storage.GetDeviceStats(deviceId)
	if errors.Is(err, &dto.ErrorNotFound{}) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "unable to load device stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dev)
}
