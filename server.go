package main

import (
	//"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Step struct {
	Type         string `json:"type"`          // 操作の種類{"heat", "add"}
	Description  string `json:"description"`   // 操作の説明(ttsで読み上げる文字列)
	Duration     string `json:"duration"`      // 操作時間[s]。加熱処理等デバイスが自動で行う操作にのみ定義されるプロパティ
	HeatStrength int64  `json:"heat_strength"` // typeがheatなときに、火力を定義する
	AddGrams     int64  `json:"add_grams"`     // typeがaddなときに、入れるものの分量をグラム指定する。
}

type Recipe struct {
	ID          int64  `json:"id"`          // レシピID(デバイスからレシピの追加を行った際は未定義でおｋ)
	UserID      int64  `json:"user_id"`     // 作成者ユーザーID(デバイスからレシピの追加を行った際は未定義でおｋ)
	DeviceID    int64  `json:"device_id"`   // デバイスID(デバイスからレシピの追加を行った際はこれを定義)
	Name        string `json:"name"`        // レシピの名前(デバイスでは未定義でおｋ）
	Description string `json:"description"` // レシピの説明(デバイスでは未定義でおｋ）
	Steps       []Step `json:"steps"`       // 操作配列
}

func main() {
	cli, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	r := gin.Default()

	// レシピアップロード
	r.POST("/recipe/upload", func(c *gin.Context) {
		var recipe Recipe
		c.BindJSON(&recipe)
		bson.UnmarshalJSON()
		fmt.Printf("[%d] %s", recipe.ID, recipe.Name)
	})

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Server Error Happened")
	}
}
