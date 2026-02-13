package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
	a.mux.HandleFunc("POST /api/v1/devices/{device_id}/heartbeat", a.heartbeatHandler)
	a.mux.HandleFunc("POST /api/v1/devices/{device_id}/stats", a.uploadDeviceStats)
	a.mux.HandleFunc("GET /api/v1/devices/{device_id}/stats", a.getDeviceStats)
	a.server.Handler = LoggingMiddleware(a.mux)
	return a
}

func (a *API) Start() error {
	return a.server.ListenAndServe()
}

func (a *API) heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	deviceId := r.PathValue("device_id")
	if deviceId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorResponse{"Device ID required"})
		return
	}
	var req dto.DeviceHeartbeat
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorResponse{"Bad JSON: " + err.Error()})
		return
	}

	heartbeatTime := req.SentAt.Unix()
	update := dto.DeviceStatUpdate{
		DeviceId:      deviceId,
		HeartbeatTime: &heartbeatTime,
	}

	err := a.storage.UpdateDeviceStats(&update)
	if errors.Is(err, &dto.ErrorNotFound{}) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}
	if err != nil {
		http.Error(w, "unable to load device stats", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) uploadDeviceStats(w http.ResponseWriter, r *http.Request) {
	deviceId := r.PathValue("device_id")
	if deviceId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorResponse{"Device ID required"})
		return
	}
	var req dto.DeviceStatsUpload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorResponse{"Bad JSON: " + err.Error()})
		return
	}

	update := dto.DeviceStatUpdate{
		DeviceId:       deviceId,
		UploadDuration: &req.UploadTime,
	}

	fmt.Println("stats: ", update)
	err := a.storage.UpdateDeviceStats(&update)
	if errors.Is(err, &dto.ErrorNotFound{}) {
		// http.Error(w, err.Error(), http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}
	if err != nil {
		http.Error(w, "unable to load device stats", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *API) getDeviceStats(w http.ResponseWriter, r *http.Request) {
	deviceId := r.PathValue("device_id")
	if deviceId == "" {
		http.Error(w, "device id required", http.StatusBadRequest)
		return
	}

	dev, err := a.storage.GetDeviceStats(deviceId)
	fmt.Println("stats: ", dev)
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
