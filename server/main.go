package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ResponseFormat struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Operation struct {
	ID           string    `json:"id"`
	LeftOperand  int       `json:"leftOperand" validate:"required,gte=-99,lte=99"`  //-- First operand
	RightOperand int       `json:"rightOperand" validate:"required,gte=-99,lte=99"` //-- Second operand
	Operator     Operator  `json:"operator" validate:"required,eq=0|eq=1|eq=2|eq=3"`
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

	response := ResponseFormat{
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

	if err := json.NewDecoder(r.Body).Decode(&newOperation); err != nil {
		response := ResponseFormat{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	validate := validator.New()
	err := validate.Struct(newOperation)

	if err != nil {
		errors := err.(validator.ValidationErrors)

		response := ResponseFormat{
			Status:  http.StatusBadRequest,
			Message: errors[0].Error(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	date := time.Now().UTC()
	newOperation.ID = uuid.New().String()
	newOperation.CreatedAt = date
	newOperation.UpdatedAt = date

	operationHistory = append(operationHistory, newOperation)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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
