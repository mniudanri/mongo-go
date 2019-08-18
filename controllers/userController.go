package controller

import (
  "mongo-go/connection"
  "net/http"
  "errors"
  "context"
  "mongo-go/models"
  "go.mongodb.org/mongo-driver/bson"
  "time"
  // "fmt"
  "encoding/json"
  "github.com/gorilla/mux"
  "go.mongodb.org/mongo-driver/bson/primitive"
)
func CreateUser(response http.ResponseWriter, request *http.Request){
  response.Header().Add("Content-Type","application/json")
  client := connection.GetClient()

  var person user.Person

  json.NewDecoder(request.Body).Decode(&person)
  err := validate(person)

  if err != nil{
    response.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(response).Encode(bson.M{"message": err.Error()})
    return
  }

  collection := client.Database("test").Collection("users")
  ctx, _ := context.WithTimeout(context.Background(), connection.GetTimeout()*time.Second)
  result, _ := collection.InsertOne(ctx, person)
  response.WriteHeader(http.StatusCreated)
  json.NewEncoder(response).Encode(result)
}

func GetUsers(response http.ResponseWriter, request *http.Request){
  response.Header().Add("Content-Type", "application/json")
  var people []user.Person
  client := connection.GetClient()
  collection := client.Database("test").Collection("users")
  ctx, _ := context.WithTimeout(context.Background(), connection.GetTimeout()*time.Second)
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

  client := connection.GetClient()
  collection := client.Database("test").Collection("users")
  ctx, _ := context.WithTimeout(context.Background(), connection.GetTimeout()*time.Second)
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

  client := connection.GetClient()

  collection := client.Database("test").Collection("users")
  ctx, _ := context.WithTimeout(context.Background(), connection.GetTimeout()*time.Second)
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
