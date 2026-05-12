package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type SensorReading struct {
	ID         string    `json:"id"`
	SensorID   string    `json:"sensor_id"`
	RecordedAt time.Time `json:"recorded_at"`
	AccX       float64   `json:"acc_x"`
	AccY       float64   `json:"acc_y"`
	AccZ       float64   `json:"acc_z"`
	PGA        *float64  `json:"pga"`
	STALTA     *float64  `json:"sta_lta"`
	IsTrigger  bool      `json:"is_trigger"`
	CreatedAt  time.Time `json:"created_at"`
}

func GetAllSensorReadings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		query := `SELECT 
		id, sensor_id, recorded_at, acc_x, acc_y, acc_z, pga, sta_lta, is_trigger, created_at
		FROM sensor_readings`

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var readings []SensorReading = make([]SensorReading, 0)

		for rows.Next() {
			var s SensorReading
			err := rows.Scan(
				&s.ID,
				&s.SensorID,
				&s.RecordedAt,
				&s.AccX,
				&s.AccY,
				&s.AccZ,
				&s.PGA,
				&s.STALTA,
				&s.IsTrigger,
				&s.CreatedAt,
			)
			if err != nil {
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			readings = append(readings, s)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, "Rows error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		duration := time.Since(start)
		fmt.Printf("Query execution time: %v\n", duration)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(readings); err != nil {
			fmt.Println("Encoding error:", err)
		}
	}
}

func GetOneSensorReadings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sensorID := r.PathValue("sensor_id")

		if sensorID == "" {
			http.Error(w, "sensor_id is required", http.StatusBadRequest)
			return
		}

		query := `SELECT 
		id, sensor_id, recorded_at, acc_x, acc_y, acc_z, pga, sta_lta, is_trigger, created_at
		FROM sensor_readings
		WHERE sensor_id = $1
		`

		rows, err := db.Query(query, sensorID)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var readings []SensorReading = make([]SensorReading, 0)

		for rows.Next() {
			var s SensorReading
			err := rows.Scan(
				&s.ID,
				&s.SensorID,
				&s.RecordedAt,
				&s.AccX,
				&s.AccY,
				&s.AccZ,
				&s.PGA,
				&s.STALTA,
				&s.IsTrigger,
				&s.CreatedAt,
			)
			if err != nil {
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			readings = append(readings, s)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, "Rows error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		duration := time.Since(start)
		fmt.Printf("Query execution time: %v\n", duration)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(readings); err != nil {
			fmt.Println("Encoding error:", err)
		}
	}
}
