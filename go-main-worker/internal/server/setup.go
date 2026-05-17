package server

import (
	"database/sql"
	"net/http"
)

func Setup(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /all-sensor-readings", GetAllSensorReadings(db))
	mux.HandleFunc("GET /one-sensor-readings/{sensor_id}", GetOneSensorReadings(db))
	mux.HandleFunc("GET /sensor-readings-join-optimized", GetOptimizedSensorReadingsJoin(db))
	mux.HandleFunc("GET /sensor-readings-join-unoptimized", GetUnoptimizedSensorReadingsJoin(db))

	return mux
}
