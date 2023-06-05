package router

import (
	"github.com/gorilla/mux"

	"contact-api/middlewares"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/contacts", middlewares.GetAllContacts).Methods("GET")

	return r;
}