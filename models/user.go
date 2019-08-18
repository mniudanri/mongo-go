package user

import (
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
  ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  Firstname string      `json:"firstname,omitempty" bson:"firstname,omitempty"`
  Lastname string       `json:"lastname,omitempty" bson:"lastname,omitempty"`
  Address string       `json:"address,omitempty" bson:"address,omitempty"`
}


// func getAll() ([]Person, errors) {
//   ctx, _ := context.WithTimeout(context.Background(), connection.GetTimeout()*time.Second)
//   curser, err := collection.Find(ctx, bson.M{})
//
// }
