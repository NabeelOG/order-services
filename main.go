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

var db *sql.DB

func initDB() {
	var err error
	
	connStr := "postgres://postgres:postgres@localhost/orderdb?sslmode=disable"

	db, err = sql.Open("postgres", connStr)
	
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
		rows, err := db.Query("SELECT id, item, quantity FROM orders ORDER BY id")
		if err != nil {
			http.Error(w, "Database error:"+ err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		
		var orders []Order
		for rows.Next() {
			var o Order
			err := rows.Scan(&o.ID, &o.Item, &o.Quantity)
			if err != nil {
				http.Error(w, "Error Scanning row: "+ err.Error(), http.StatusInternalServerError)
				return
			}
			orders = append(orders, o)
		}

		json.NewEncoder(w).Encode(orders)
	case "POST":
		var newOrder Order
		err := json.NewDecoder(r.Body).Decode(&newOrder)
		if err!= nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if newOrder.Item=="" {
			http.Error(w, "Item cant be empty", http.StatusBadRequest)
			return
		}
		if newOrder.Quantity<=0 {
			http.Error(w, "Quantity must be positive", http.StatusBadRequest)
			return
		}

		var id int
		err = db.QueryRow(
			"INSERT INTO orders (item, quantity) VALUES ($1, $2) RETURNING id", newOrder.Item, newOrder.Quantity,
		).Scan(&id)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		newOrder.ID = id
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

	switch r.Method {
	case "GET":
		//SINGLE ORDER FROM DATABASE
		var order Order
		err = db.QueryRow(
			"SELECT id, item, quantity FROM orders WHERE id = $1", id,
		).Scan(&order.ID, &order.Item, &order.Quantity)

		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database error"+err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(order)

	case "PUT":
		var updateOrder Order
		err := json.NewDecoder(r.Body).Decode(&updateOrder)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if updateOrder.Item == "" {
			http.Error(w, "Invalid Request Body", http.StatusBadRequest)
			return
		}

		if updateOrder.Quantity <= 0 {
			http.Error(w, "Quantity must be positive", http.StatusBadRequest)
			return
		}

		result, err := db.Exec(
			"UPDATE orders SET item = $1, quantity = $2 WHERE id = $3",
			updateOrder.Item, updateOrder.Quantity, id,
		)

		if err!= nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		updateOrder.ID = id
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updateOrder)

	case "DELETE":
		result, err := db.Exec("DELETE FROM orders WHERE id = $1", id)
		if err != nil {
			http.Error(w, "Database Error"+ err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}	
}