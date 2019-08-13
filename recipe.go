package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

func recipeUpload(c *gin.Context) {
  var recipe Recipe
  buf := make([]byte, 2048)
  c.Request.Body.Read(buf)
  bson.UnmarshalExtJSON(buf, false, recipe)
  if _, err := recipeDB.InsertOne(context.TODO(), recipe); err != nil{
	c.Status(500)
  }
  c.Status(200)
}

func recipeGet(c *gin.Context) {
  ctx := context.Background()
  var recipe bson.M
  filter := bson.M{"_id": c.Param("id")}
  if err := userDB.FindOne(ctx, filter).Decode(&recipe); err != nil{
	c.Status(500)
  }
  c.JSON(200, recipe)
}
