package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// APIResponse adalah struktur standar untuk respons JSON
type APIResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func main() {
	// 1. Inisialisasi ServeMux (Multiplexer / Router bawaan Go)
	router := http.NewServeMux()

	// 2. Registrasi Endpoint (Menggunakan fitur routing modern Go 1.22+)
	router.HandleFunc("GET /", handleHome)
	router.HandleFunc("GET /health", handleHealth)
	router.HandleFunc("GET /users/{id}", handleGetUser) // Mengambil path parameter {id}
	router.HandleFunc("POST /echo", handleEcho)          // Endpoint untuk menerima JSON

	// 3. Konfigurasi HTTP Server
	// Kami mendefinisikan struct http.Server secara eksplisit untuk kontrol timeout yang lebih baik.
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("Server running on http://localhost%s\n", server.Addr)
	
	// Jalankan server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
	}
}

// ==========================================
// HANDLERS
// ==========================================

// handleHome menangani request ke root path "/"
func handleHome(w http.ResponseWriter, r *http.Request) {
	// Menolak request jika path tidak cocok persis dengan "/" (mencegah catch-all default)
	if r.URL.Path != "/" {
		respondWithError(w, http.StatusNotFound, "Page not found")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Message: "Welcome to the lightweight, dependency-free Go service!",
		Status:  http.StatusOK,
	})
}

// handleHealth mengembalikan status kesehatan aplikasi (useful untuk Docker healthcheck)
func handleHealth(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "UP"})
}

// handleGetUser mengekstrak path parameter {id} tanpa library tambahan
func handleGetUser(w http.ResponseWriter, r *http.Request) {
	// Fitur Go 1.22+: r.PathValue() mengekstrak parameter dari URL pattern
	userID := r.PathValue("id")
	if userID == "" {
		respondWithError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Message: fmt.Sprintf("Successfully fetched data for User ID: %s", userID),
		Status:  http.StatusOK,
	})
}

// handleEcho membaca JSON request dan mengembalikannya kembali ke client
func handleEcho(w http.ResponseWriter, r *http.Request) {
	var input map[string]interface{}

	// Menghindari memory leak dengan menggunakan json.NewDecoder langsung pada r.Body
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	defer r.Body.Close()

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "echoed",
		"data":   input,
	})
}

// ==========================================
// HELPERS (UTILITAS ENKODING JSON)
// ==========================================

// respondWithJSON mengirimkan payload dalam format JSON dengan status code yang tepat
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	// Menggunakan json.NewEncoder langsung ke Writer (lebih cepat dibanding json.Marshal)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

// respondWithError mengirimkan standarisasi pesan error JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, APIResponse{
		Message: message,
		Status:  code,
	})
}
