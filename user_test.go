package main

import (
	"context"
	"testing"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestUser(t *testing.T) {
  t.Logf("Hi")
  ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
  defer cancel()
  cli, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
  if err != nil {
	t.Fatal(err)
  }
  userDB := cli.Database("mesistant").Collection("user")
  ctx = context.Background()
  filter := bson.M{}
  cur, err := userDB.Find(ctx, filter)
  if err != nil {
	t.Error(err)
  }
  defer cur.Close(ctx)
  for cur.Next(ctx) {
	var result bson.M
	err := cur.Decode(&result)
	if err != nil {
	  t.Error(err)
	}
	t.Log(result)
  }
}

