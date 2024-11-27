package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/google/uuid"  
	"math"
	"net/http"
	"os"  
	"strconv"
	"strings"
	"time"
)

type Receipt struct {
	Retailer     string    `json:"retailer"`
	PurchaseDate string    `json:"purchaseDate"`
	PurchaseTime string    `json:"purchaseTime"`
	Items        []Item    `json:"items"`
	Total        string    `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var ReceiptStore = make(map[string]int)

func main() {
	r := mux.NewRouter()

	
	r.HandleFunc("/receipts/process", ProcessReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", GetPoints).Methods("GET")

	
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	fmt.Printf("Server is running on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, r)
}

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt

	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	receiptID := uuid.New().String() // Generates a new unique ID
	points := calculatePoints(receipt)
	ReceiptStore[receiptID] = points

	response := map[string]string{"id": receiptID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receiptID := vars["id"]

	points, found := ReceiptStore[receiptID]
	if !found {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	response := map[string]int{"points": points}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}


func calculatePoints(receipt Receipt) int {
	points := 0
	total, _ := strconv.ParseFloat(receipt.Total, 64)

	points += countAlphanumeric(receipt.Retailer)

	if total == math.Floor(total) {
		points += 50
	}

	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		itemPrice, _ := strconv.ParseFloat(item.Price, 64)
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points += int(math.Ceil(itemPrice * 0.2))
		}
	}

	if isOddDay(receipt.PurchaseDate) {
		points += 6
	}

	if isBetween2And4PM(receipt.PurchaseTime) {
		points += 10
	}

	return points
}

func countAlphanumeric(input string) int {
	count := 0
	for _, char := range input {
		if isAlphanumeric(char) {
			count++
		}
	}
	return count
}

func isAlphanumeric(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func isOddDay(dateStr string) bool {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}
	return date.Day()%2 != 0
}

func isBetween2And4PM(timeStr string) bool {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return false
	}
	return t.Hour() >= 14 && t.Hour() < 16
}
