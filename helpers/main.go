package main

import (
	"fmt"
	"log"
	"net/http"

	"helpers/filter_values"
)

func main() {
	http.HandleFunc("/getFilterValues", filter_values.GetFilterValues)
	fmt.Println("Server is listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}