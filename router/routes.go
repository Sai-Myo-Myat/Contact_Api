package router

import (
	"github.com/gorilla/mux"

	"contact-api/middlewares"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/contacts", middlewares.GetAllContacts).Methods("GET")
	r.HandleFunc("/api/contacts", middlewares.CreateContact).Methods("POST")
	r.HandleFunc("/api/contacts/{id}", middlewares.GetContact).Methods("GET")

	return r;
}