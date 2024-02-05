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
	Id           string    `json:"id"`
	LeftOperand  int       `json:"left_operand"`  //-- First operand
	RightOperand int       `json:"right_operand"` //-- Second operand
	Operator     Operator  `json:"operator"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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

func main() {
	const PORT string = "3000"
	router := mux.NewRouter().StrictSlash(true)

	now := time.Now().UTC()
	operationHistory = append(operationHistory,
		Operation{
			Id:           uuid.New().String(),
			LeftOperand:  22,
			RightOperand: 2,
			Operator:     ADD,
			CreatedAt:    now,
			UpdatedAt:    now,
		})

	router.HandleFunc("/health", healthCheck).Methods("GET")
	router.HandleFunc("/history", getOperations).Methods("GET")

	fmt.Println("Server running on port:", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
