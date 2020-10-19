package main

import (
	"log"
	"net/http"

	"gopkg.in/go-playground/validator.v9"

	"github.com/gorilla/mux"
	"github.com/newestuser/faceit/api"
)

func main() {
	validate := validator.New()

	router := mux.NewRouter()
	router.Handle("/users", api.UserRegHandler(validate, nil)).Methods("POST")
	router.Handle("/users/{id}", api.UserGetHandler(nil)).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
