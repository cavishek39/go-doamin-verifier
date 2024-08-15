package main

import (
	"fmt"
	"go-postgres-stock-api/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router() // call the router
	fmt.Println("Starting server on the port 8080...")

	log.Fatal(http.ListenAndServe(":8080", r)) // Run the server
}