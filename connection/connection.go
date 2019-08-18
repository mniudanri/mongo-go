package connection

import (
  "context"
  "time"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

var (
  client *mongo.Client
  timeout_max time.Duration = 10
)

func Connect() (*mongo.Client, error){
  ctx, _ := context.WithTimeout(context.Background(), timeout_max*time.Second)
  client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

  return client, nil
}

func GetClient()(*mongo.Client){
  return client
}

func GetTimeout() time.Duration {
  return timeout_max
}
