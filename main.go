package main

import (
	"contact-api/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	fmt.Println("Server is running on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}