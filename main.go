package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "go.mongodb.org/mongo-driver/mongo"
  "mongo-go/connection"
  "mongo-go/controllers"
)
var client *mongo.Client

func main() {
  client,_ = connection.Connect()

  router := mux.NewRouter()
  router.HandleFunc("/person", controller.CreateUser).Methods("POST")
  router.HandleFunc("/person", controller.GetUsers).Methods("GET")
  router.HandleFunc("/person/{id}", controller.GetUserById).Methods("GET")
  router.HandleFunc("/person/{id}", controller.DeleteUserById).Methods("DELETE")

  http.ListenAndServe(":12345", router)
}
