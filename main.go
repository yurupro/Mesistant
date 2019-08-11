package main

import (
	"fmt"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Step struct {
	Type         string `bson:"type"`          // 操作の種類{"heat", "add"}
	Description  string `bson:"description"`   // 操作の説明(ttsで読み上げる文字列)
	Duration     string `bson:"duration"`      // 操作時間[s]。加熱処理等デバイスが自動で行う操作にのみ定義されるプロパティ
	HeatStrength int64  `bson:"heat_strength"` // typeがheatなときに、火力を定義する
	AddGrams     int64  `bson:"add_grams"`     // typeがaddなときに、入れるものの分量をグラム指定する。
}

type Recipe struct {
	ID          int64  `bson:"id"`          // レシピID(デバイスからレシピの追加を行った際は未定義でおｋ)
	UserID      int64  `bson:"user_id"`     // 作成者ユーザーID(デバイスからレシピの追加を行った際は未定義でおｋ)
	DeviceID    int64  `bson:"device_id"`   // デバイスID(デバイスからレシピの追加を行った際はこれを定義)
	Name        string `bson:"name"`        // レシピの名前(デバイスでは未定義でおｋ）
	Description string `bson:"description"` // レシピの説明(デバイスでは未定義でおｋ）
	Steps       []Step `bson:"steps"`       // 操作配列
}

func initDB() (*mongo.Client, error) {
  cli, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
  if err != nil {
	return nil, err
  }
  return cli, nil
}

func insertRecipe(col *mongo.Collection, recipe Recipe) error {
  if _, err := col.InsertOne(context.TODO(), recipe); err != nil {
	return err
  }
  return nil
}

func main() {
	r := gin.Default()
	dbcli, err := initDB()
	if err != nil {
	  fmt.Println(err)
	  return
	}
	col := dbcli.Database("test").Collection("trailers")

	// レシピアップロード
	r.POST("/recipe/upload", func(c *gin.Context) {
		var recipe Recipe
		buf := make([]byte, 2048)
		c.Request.Body.Read(buf)
		bson.UnmarshalExtJSON(buf, false, recipe)
		if err := insertRecipe(col, recipe); err != nil {
		}

		fmt.Printf("[%d] %s", recipe.ID, recipe.Name)
	})

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Server Error Happened")
	}
}
