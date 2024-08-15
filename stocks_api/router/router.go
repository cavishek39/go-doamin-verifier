package router

import (
	"github.com/gorilla/mux"
	"go-postgres-stock-api/middleware"
)

func Router () * mux.Router{
	router := mux.NewRouter() // create a new router

	router.HandleFunc("/api/stocks", middleware.GetStocks).Methods("GET", "OPTIONS") // get all stocks
	// router.HandleFunc("/api/stocks/{id}", middleware.GetStock).Methods("GET", "OPTIONS") // get a stock by id
	router.HandleFunc("/api/new-stock", middleware.CreateStock).Methods("POST", "OPTIONS") // create a new stock
	// router.HandleFunc("/api/update-stock/{id}", middleware.UpdateStock).Methods("PUT", "OPTIONS") // update a stock by id
	router.HandleFunc("/api/delete-stock/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS") // delete a stock by id
	router.HandleFunc("/api/delete-all", middleware.DeleteAllStocks).Methods("DELETE", "OPTIONS") // delete all stocks

	return router
}