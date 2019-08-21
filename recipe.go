package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-session/gin-session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JSONMultiRecipe struct {
	Array []Recipe `json:"array"`
}
type Step struct {
	Type         string `bson:"type" json:"type"`                   // 操作の種類{"heat", "add"}
	Description  string `bson:"description" json:"description"`     // 操作の説明(ttsで読み上げる文字列)
	Duration     int64  `bson:"duration" json:"duration"`           // 操作時間[s]。加熱処理等デバイスが自動で行う操作にのみ定義されるプロパティ
	HeatStrength int64  `bson:"heat_strength" json:"heat_strength"` // typeがheatなときに、火力を定義する
	AddGrams     int64  `bson:"add_grams" json:"add_grams"`         // typeがaddなときに、入れるものの分量をグラム指定する。
}

type Recipe struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"` // レシピID(デバイスからレシピの追加を行った際は未定義でおｋ)
	UserID      string             `bson:"user_id" json:"user_id"`
	Name        string             `bson:"name" json:"name"`               // レシピの名前(デバイスでは未定義でおｋ）
	Description string             `bson:"description" json:"description"` // レシピの説明(デバイスでは未定義でおｋ）
	Steps       []Step             `bson:"steps" json:"steps"`             // 操作配列
}

func recipeUpload(c *gin.Context) {
	var recipe Recipe
	recipe.ID = primitive.NewObjectID()
	if err := c.BindJSON(&recipe); err != nil {
		c.Status(500)
		return
	}
	fmt.Println(recipe)
	if _, err := recipeDB.InsertOne(context.TODO(), recipe); err != nil {
		c.Status(500)
		return
	}
	c.JSON(200, recipe)
}

func recipeGet(c *gin.Context) {
	ctx := context.Background()
	var recipe Recipe
	recipe.ID = primitive.NilObjectID
	recipeID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.Status(400)
		return
	}
	filter := bson.M{"_id": recipeID}
	if err := recipeDB.FindOne(ctx, filter).Decode(&recipe); err != nil {
		fmt.Println(err)
		c.Status(500)
		return
	}
	c.JSON(200, recipe)
}

func recipeGetByUser(c *gin.Context) {
	ctx := context.Background()
	recipes := make([]Recipe, 0, 20)
	userID := c.Param("id")
	fmt.Println(userID)
	filter := bson.M{"user_id": userID}
	cur, err := recipeDB.Find(ctx, filter)
	if err != nil {
		c.Status(500)
		return
	}

	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		fmt.Println(recipe)
		recipes = append(recipes, recipe)
	}
	c.JSON(200, JSONMultiRecipe{Array: recipes})
}
func recipeAddQueue(c *gin.Context) {
	ctx := context.TODO()
	store := ginsession.FromContext(c)
	userID, ok := store.Get("user")
	if !ok {
		c.Status(403)
		return
	}
	var recipe Recipe
	var device Device
	recipeID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := recipeDB.FindOne(ctx, bson.M{"_id": recipeID}).Decode(&recipe); err != nil {
		c.Status(404)
		return
	}
	if err := deviceDB.FindOne(ctx, bson.M{"user_id": userID.(string)}).Decode(&device); err != nil {
		c.Status(400)
		return
	}
	device.Recipe = recipe
	c.JSON(200, recipe)
}
