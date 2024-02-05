package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HealthCheckResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthCheckResponse{
		Status:  http.StatusOK,
		Message: "OK",
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	const PORT string = "3000"
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/health", healthCheck)

	fmt.Println("Server running on port:", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
