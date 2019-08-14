package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-session/gin-session"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var recipeDB *mongo.Collection
var userDB *mongo.Collection
var deviceDB *mongo.Collection

func initDB() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func main() {
	r := gin.Default()
	r.Use(ginsession.New())
	r.Use(static.Serve("/", static.LocalFile("web", false)))

	mongodb, err := initDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	recipeDB = mongodb.Database("mesistant").Collection("recipe")
	userDB = mongodb.Database("mesistant").Collection("user")
	deviceDB = mongodb.Database("mesistant").Collection("device")

	// レシピアップロード
	r.POST("/recipe/upload", recipeUpload)
	r.GET("/recipe/:id", recipeGet)
	r.GET("/user/:id/recipes", recipeGetByUser)

	// ユーザー追加
	r.POST("/user/add", userAdd)
	r.POST("/user/login", userLogin) //{"user": "User id", "password": "password"}(SSLだからボディーにJSON載せよう)
	r.POST("/user/logout", userLogout)

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Server Error Happened")
	}
}
