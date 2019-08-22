package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"time"
)

type Config struct {
	DatabaseURL string `json:"database_url"`
}

var recipeDB *mongo.Collection
var userDB *mongo.Collection
var deviceDB *mongo.Collection

func initDB(url string) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(url))

	if err != nil {
		return nil, err
	}

	if err = cli.Ping(context.TODO(), nil); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return cli, nil
}

func main() {
	var config Config
	// 設定ファイル読み込み
	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
		config.DatabaseURL = "mongodb://mesistant_db:27017"
	}

	if err := json.Unmarshal(raw, &config); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config)

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.Use(static.Serve("/", static.LocalFile("web", false)))

	mongodb, err := initDB(config.DatabaseURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipeDB = mongodb.Database("mesistant").Collection("recipe")
	userDB = mongodb.Database("mesistant").Collection("user")
	deviceDB = mongodb.Database("mesistant").Collection("device")

	// レシピアップロード
	r.POST("/recipe", recipeUpload)
	// IDからレシピ取得
	r.GET("/recipe/:id", recipeGet)
	// キューにレシピを追加
	r.POST("/recipe/:id/add_queue", recipeAddQueue)
	// ユーザーIDからそのユーザーのレシピを取得
	r.GET("/user/:id/recipes", recipeGetByUser)

	// ユーザー追加
	r.POST("/user/add", userAdd)
	// ユーザーログイン
	r.POST("/user/login", userLogin) //{"user": "User id", "password": "password"}(SSLだからボディーにJSON載せよう)
	// ユーザーログアウト
	r.POST("/user/logout", userLogout)

	// デバイスの登録
	r.POST("/device/register", registerDevice)
	// デバイスからのキューの確認
	r.GET("/device/queue/:id", getDeviceQueue)

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Server Error Happened")
	}
}
