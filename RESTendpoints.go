package main

import (
	"encoding/json"
	"log"
	"net/http"

	"workout_ledger/database"
	"workout_ledger/services"
)

// MessageResponse represents the response structure for the helloHandler
type MessageResponse struct {
	Message string `json:"message"`
}

// HealthResponse represents the response structure for the healthHandler
type HealthResponse struct {
	Status string `json:"status"`
}

// helloHandler godoc
// @Summary      Hello endpoint
// @Description  Returns a simple hello message.
// @Tags         demo
// @Produce      json
// @Success      200  {object}  MessageResponse
// @Failure      405  {string}  string  "method not allowed"
// @Router       /hello [get]
func helloHandler(writerHttp http.ResponseWriter, requestHttp *http.Request) {
	if requestHttp.Method != http.MethodGet {
		writerHttp.Header().Set("Allow", http.MethodGet)
		http.Error(writerHttp, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writerHttp.Header().Set("Content-Type", "application/json")
	response := MessageResponse{Message: "Hello, world!"}
	if err := json.NewEncoder(writerHttp).Encode(response); err != nil {
		log.Printf("failed to encode json: %v", err)
		http.Error(writerHttp, "internal server error", http.StatusInternalServerError)
	}
}

// healthHandler godoc
// @Summary      Health check
// @Description  Returns ok if the service is running.
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Failure      405  {string}  string  "method not allowed"
// @Router       /health [get]
func healthHandler(writerHttp http.ResponseWriter, requestHttp *http.Request) {
	if requestHttp.Method != http.MethodGet {
		writerHttp.Header().Set("Allow", http.MethodGet)
		http.Error(writerHttp, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writerHttp.Header().Set("Content-Type", "application/json")
	response := HealthResponse{Status: "ok"}
	if err := json.NewEncoder(writerHttp).Encode(response); err != nil {
		log.Printf("failed to encode json: %v", err)
		http.Error(writerHttp, "internal server error", http.StatusInternalServerError)
	}
}

// testDBselect godoc
// @Summary      database connection and select test
// @Description  Returns ok if the service is running.
// @Tags         health
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Failure      405  {string}  string  "method not allowed"
// @Router       /dbtest [get]
func testDBselect(writerHttp http.ResponseWriter, requestHttp *http.Request) {
	if requestHttp.Method != http.MethodGet {
		writerHttp.Header().Set("Allow", http.MethodGet)
		http.Error(writerHttp, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writerHttp.Header().Set("Content-Type", "application/json")

	svc := services.TestDBService{DB: database.DB()}
	if err := svc.SelectAllTestRows(); err != nil {
		log.Printf("db test failed: %v", err)
		http.Error(writerHttp, "database query failed", http.StatusInternalServerError)
		return
	}

	response := HealthResponse{Status: "ok"}
	if err := json.NewEncoder(writerHttp).Encode(response); err != nil {
		log.Printf("failed to encode json: %v", err)
		http.Error(writerHttp, "internal server error", http.StatusInternalServerError)
	}
}
