package main

import (
  "net/http"
  "github.com/gorilla/mux"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/bson"
  "context"
  "go.mongodb.org/mongo-driver/mongo/options"
  "time"
  "errors"
  "encoding/json"
  "mongo-go/models"
)

var client *mongo.Client

func CreateUser(response http.ResponseWriter, request *http.Request){
  response.Header().Add("Content-Type","application/json")
  var person user.Person

  json.NewDecoder(request.Body).Decode(&person)
  err := validate(person)

  if err != nil{
    response.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(response).Encode(bson.M{"message": err.Error()})
    return
  }

  collection := client.Database("test").Collection("users")
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  result, _ := collection.InsertOne(ctx, person)
  response.WriteHeader(http.StatusCreated)
  json.NewEncoder(response).Encode(result)
}

func GetUsers(response http.ResponseWriter, request *http.Request){
  response.Header().Add("Content-Type", "application/json")
  var people []user.Person
  collection := client.Database("test").Collection("users")
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  curser, err := collection.Find(ctx, bson.M{})

  if(err != nil){
    response.WriteHeader(http.StatusInternalServerError)
    response.Write([]byte(`{"message":"`+ err.Error() +`"}`))
    return
  }

  defer curser.Close(ctx)

  for curser.Next(ctx){
    var person user.Person
    curser.Decode(&person)
    people = append(people, person)
  }

  if err := curser.Err(); err != nil {
    response.WriteHeader(http.StatusInternalServerError)
    response.Write([]byte(`{"message":"`+ err.Error() +`"}`))
    return
  }
  json.NewEncoder(response).Encode(people)
}

func GetUserById(response http.ResponseWriter, request *http.Request){
  response.Header().Add("Content-Type", "application/json")
  params := mux.Vars(request)
  id, _ := primitive.ObjectIDFromHex(params["id"])
  var person user.Person
  collection := client.Database("test").Collection("users")
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  err := collection.FindOne(ctx, user.Person{ID: id}).Decode(&person)

  if err != nil {
    response.WriteHeader(http.StatusInternalServerError)
    response.Write([]byte(`{"message":"`+ err.Error() +`"}`))
    return
  }

  json.NewEncoder(response).Encode(person)
}

func DeleteUserById(response http.ResponseWriter, request *http.Request) {
  response.Header().Add("Content-Type", "application/json")
  params := mux.Vars(request)
  id, _ := primitive.ObjectIDFromHex(params["id"])
  collection := client.Database("test").Collection("users")
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  res, err := collection.DeleteOne(ctx, bson.M{"_id": id})

  if err != nil {
    response.WriteHeader(http.StatusInternalServerError)
    response.Write([]byte(`{"message":"`+ err.Error() +`"}`))
    return
	}

  if res.DeletedCount == 0 {
    response.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(response).Encode(bson.M{"message":"Id not found!"})
    return
  }

  response.WriteHeader(http.StatusOK)
  json.NewEncoder(response).Encode(bson.M{"message": "Delete success", "_id": res.DeletedCount})
  return
}

func validate(data user.Person) error{
  if data.Firstname == "" {
    return errors.New("firstname is required")
  }else if data.Lastname == "" {
    return errors.New("lastname is required")
  }

  return nil
}


func main() {
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

  router := mux.NewRouter()
  router.HandleFunc("/person", CreateUser).Methods("POST")
  router.HandleFunc("/person", GetUsers).Methods("GET")
  router.HandleFunc("/person/{id}", GetUserById).Methods("GET")
  router.HandleFunc("/person/{id}", DeleteUserById).Methods("DELETE")

  http.ListenAndServe(":12345", router)
}
