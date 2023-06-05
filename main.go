package main

import (
	"contact-api/router"
	"contact-api/middlewares"
	"fmt"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	middlewares.CreateDB()
	fmt.Println("Server is running on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}