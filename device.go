package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	//"github.com/go-session/gin-session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Device struct {
  ID	primitive.ObjectID	`bson:"_id" json:"id"`
  UserID	string	`bson:"user_id" json:"user_id"`
  Recipe	Recipe `bson:"recipe" json:"recipe"`
}


func registerDevice(c *gin.Context) { // POST by device
  var device Device
  device.ID = primitive.NewObjectID()
  if err := c.BindJSON(&device); err != nil {
	c.Status(400)
	fmt.Println(err)
	return
  }
  ctx := context.TODO()
  res, err := deviceDB.InsertOne(ctx, device)
  if err != nil {
	c.Status(400)
	fmt.Println(err)
	return
  }
  device.ID = res.InsertedID.(primitive.ObjectID)
  c.JSON(200, device)
}

func getDeviceQueue(c *gin.Context) { // GET by device
  deviceID, err := primitive.ObjectIDFromHex(c.Param("id"))
  if err != nil {
	c.Status(400)
	fmt.Println(err)
	return
  }
  ctx := context.TODO()
  var device Device
  if err := deviceDB.FindOne(ctx, bson.M{"_id": deviceID}).Decode(&device); err != nil {
	c.Status(500)
	fmt.Println(err)
	return
  }
  var recipe Recipe
  recipe = device.Recipe
  c.JSON(200, recipe)
}