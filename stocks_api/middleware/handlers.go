package middleware

import (
	"database/sql"
	"encoding/json"
	"go-postgres-stock-api/models"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"fmt"
    "strconv"
    "github.com/gorilla/mux"

)

// Type of the response message to be sent back to the client
type response struct {
	ID int64 `json:"id,omitempty"`
	Status int `json:"status"`
	Message string `json:"message"`
    Data interface{} `json:"data,omitempty"`
}

func CreateConnection() *sql.DB {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    connStr := os.Getenv("POSTGRES_URL")
	// fmt.Println(connStr)
	
    db, err := sql.Open("postgres", connStr)

    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }

    return db
}

// CreateStock creates a new stock
func CreateStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }


    var stock models.Stock
    var res response


    err := json.NewDecoder(r.Body).Decode(&stock)
    if err != nil {
        res = response{
            Status:  http.StatusBadRequest,
            Message: "Invalid request payload",
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(res)
        return
    }

    stock.ID = 0

    pool := CreateConnection()
    defer pool.Close()

    err = pool.QueryRow("INSERT INTO stocks (name, price, company) VALUES ($1, $2, $3) RETURNING id;", stock.Name, stock.Price, stock.Company).Scan(&stock.ID)
    if err != nil {
        res = response{
            Status:  http.StatusBadRequest,
            Message: "Server error",
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(res)
        return
    }

    res = response{
        ID:      stock.ID,
        Status:  http.StatusCreated,
        Message: "Stock created successfully",
        Data:    stock,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(res)
}

// GetStocks gets all stocks
func GetStocks(	w http.ResponseWriter, r *http.Request) {
	pool := CreateConnection()
	defer pool.Close()

	rows, err := pool.Query("SELECT * FROM stocks")

	if err != nil {
		panic(err)
	}

	stocks := []models.Stock{}

	for rows.Next(){ 
		var stock models.Stock
		err = rows.Scan(&stock.ID, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			panic(err)
		}

		stocks = append(stocks, stock)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(stocks); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// Delete a stock
func DeleteStock(w http.ResponseWriter, r *http.Request) {
    // Check if the method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

    // Create a connection to the database
    pool := CreateConnection()
    defer pool.Close()

    vars := mux.Vars(r)
    // Get the id from the request
    idStr := vars["id"]
    // fmt.Println("idStr => ", idStr)

    if idStr == "" {
        http.Error(w, "ID parameter is missing", http.StatusBadRequest)
        return
    }

    // Convert the id to an integer
    id, errAtoi := strconv.Atoi(idStr)
    
    if errAtoi != nil {
        http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
        return
    }

    // Execute the delete query
    _, err := pool.Exec("DELETE FROM stocks WHERE id = $1", id)

    if err != nil {
        panic(err)
    }

    // Send a response to the client
    w.Header().Set("Content-type", "application/json")

    res := response{
        Status:  http.StatusOK,
        Message: "Stock deleted successfully",
    }

    fmt.Println("Stock deleted successfully")

    if err := json.NewEncoder(w).Encode(res); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// Delete all stocks
func DeleteAllStocks(w http.ResponseWriter, r *http.Request) {
    // Check if the method is DELETE
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Create a connection to the database
    pool := CreateConnection()
    defer pool.Close() // Close the connection when the function returns

    // Execute the delete query
    _, err := pool.Exec("DELETE FROM stocks")

    if err != nil {
        panic(err)
        http.Error(w, "Server error", http.StatusInternalServerError)
    }

    // Send a response to the client
    w.Header().Set("Content-type", "application/json")

    res := response {
        Status: http.StatusOK,
        Message: "All stocks deleted successfully",
    }

    if err := json.NewEncoder(w).Encode(res); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    fmt.Println("All stocks deleted successfully")
}