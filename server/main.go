package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HealthCheckResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Operation struct {
	ID           string    `json:"id"`
	LeftOperand  int       `json:"leftOperand"`  //-- First operand
	RightOperand int       `json:"rightOperand"` //-- Second operand
	Operator     Operator  `json:"operator"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Operator int

const (
	ADD Operator = iota
	SUBT
	MULT
	DIV
)

var operationHistory = make([]Operation, 0)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthCheckResponse{
		Status:  http.StatusOK,
		Message: "OK",
	}

	json.NewEncoder(w).Encode(response)
}

func getOperations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(operationHistory)
}

func createOperation(w http.ResponseWriter, r *http.Request) {
	var newOperation Operation

	error := json.NewDecoder(r.Body).Decode(&newOperation)

	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	date := time.Now().UTC()
	newOperation.ID = uuid.New().String()
	newOperation.CreatedAt = date
	newOperation.UpdatedAt = date

	operationHistory = append(operationHistory, newOperation)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newOperation)
}

func main() {
	const PORT string = "3000"
	router := mux.NewRouter().StrictSlash(true)

	now := time.Now().UTC()
	//-- Hardcoded data --
	operationHistory = append(operationHistory,
		Operation{
			ID:           uuid.New().String(),
			LeftOperand:  22,
			RightOperand: 2,
			Operator:     ADD,
			CreatedAt:    now,
			UpdatedAt:    now,
		})

	//-- Routes --
	router.HandleFunc("/health", healthCheck).Methods("GET")
	router.HandleFunc("/history", getOperations).Methods("GET")
	router.HandleFunc("/history", createOperation).Methods("POST")

	fmt.Println("Server running on port:", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
