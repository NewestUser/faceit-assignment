package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/newestuser/faceit/api"
	"github.com/newestuser/faceit/kfka"
	"github.com/newestuser/faceit/mgo"
	"gopkg.in/go-playground/validator.v9"
)

func main() {

	appPort := env("PORT", "8080")
	mgoHost := env("MONGO_HOST", "127.0.0.1")
	mgoPort, _ := strconv.Atoi(env("MONGO_PORT", "27017"))
	mgoUser := env("MONGO_USER", "root")
	mgoPswd := env("MONGO_PASSWORD", "example")
	mgoDB := env("MONGO_DB_NAME", "faceit")
	kfkaHost := env("KAFKA_HOST", "127.0.0.1")

	db := mgo.NewDb(mgoUser, mgoPswd, mgoHost, mgoPort, mgoDB)
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect()

	producer, err := kfka.NewKafkaProducer(kfkaHost)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	userEmitter := kfka.NewKafkaUserEventEmitter(producer)

	repo := mgo.NewUserRepository(db, userEmitter)
	validate := validator.New()

	router := mux.NewRouter()
	router.Handle("/users", api.UserRegHandler(validate, repo)).Methods("POST")
	router.Handle("/users/{id}", api.UserGetHandler(repo)).Methods("GET")
	router.Handle("/users/{id}", api.UserUpdateHandler(validate, repo)).Methods("PUT")

	fmt.Printf("Starting FACEIT User API on :%s\n", appPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appPort), router))
}

func env(key string, fallback string) string {
	if envVar := os.Getenv(key); envVar != "" {
		return envVar
	}
	return fallback
}
