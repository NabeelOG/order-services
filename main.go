package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Order struct {
	ID       int    `json:"id"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

var orders = []Order{
	{ID: 1, Item: "laptop", Quantity: 1},
	{ID: 2, Item: "tablet", Quantity: 2},
}

var db *sql.DB

func initDB() {
	connStr := "postgres://postgres:postgres@localhost/orderdb?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Connected to PostgresSQL database!")
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/orders", ordersHandler)
	http.HandleFunc("/orders/", orderHandler)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(orders)
	case "POST":
		var newOrder Order
		err := json.NewDecoder(r.Body).Decode(&newOrder)
		if err!= nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		newOrder.ID = len(orders) + 1
		orders = append(orders, newOrder)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newOrder)
	default:
		http.Error(w,"Method not allowed", http.StatusMethodNotAllowed)
	}
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	index := -1
	for _,order := range orders {
		if order.ID == id {
			index = id
			break
		}
	}

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(orders[index])

	case "PUT":
		var updateOrder Order
		err := json.NewDecoder(r.Body).Decode(&updateOrder)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		updateOrder.ID = id
		orders[index] = updateOrder

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updateOrder)

	case "DELETE":
		orders = append(orders[:index], orders[index+1:]...)
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}	
}