package main

import (
	"github.com/gorilla/mux"
	"github.com/newestuser/faceit/api"
	"github.com/newestuser/faceit/mgo"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
)

func main() {
	db := mgo.NewDb("root", "example", "127.0.0.1", 27017, "faceit")
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}

	repo := mgo.NewUserRepository(db)
	validate := validator.New()

	router := mux.NewRouter()
	router.Handle("/users", api.UserRegHandler(validate, repo)).Methods("POST")
	router.Handle("/users/{id}", api.UserGetHandler(repo)).Methods("GET")
	router.Handle("/users/{id}", api.UserUpdateHandler(validate, repo)).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}
