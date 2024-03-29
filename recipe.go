package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	  fmt.Println(err)
		c.Status(400)
		return
	}
	fmt.Println(recipe)
	ctx := context.TODO()
	var user User
	userID, err := primitive.ObjectIDFromHex(recipe.UserID)
	if err != nil {
	  c.Status(400)
	  return
	}
	if err := userDB.FindOne(ctx, bson.M{"_id": userID}).Decode(&user); err != nil {
	  c.Status(500)
	  return
	}

	if _, err := recipeDB.InsertOne(context.TODO(), recipe); err != nil {
	  fmt.Println(err)
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

func recipeDelete(c *gin.Context) {
  session := sessions.Default(c)
  userID := session.Get("user")
  if userID == nil {
	c.Status(403)
	return
  }
  ctx := context.TODO()
  recipeID, _:= primitive.ObjectIDFromHex(c.Param("id"))

  if _, err := recipeDB.DeleteOne(ctx, bson.M{"_id": recipeID, "user_id": userID}); err != nil{
	c.Status(400)
	return
  }
  c.Status(200)
}

func recipeAll(c *gin.Context) {
	ctx := context.Background()
	recipes := make([]Recipe, 0, 20)

    options := options.Find()
    options.SetSort(bson.D{{"_id", -1}})
    
	cur, err := recipeDB.Find(ctx, bson.M{}, options)
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

func recipeUpdate(c *gin.Context) {
    session := sessions.Default(c)
    userID := session.Get("user")
    if userID == nil {
        c.Status(403)
        return
    }
    ctx := context.TODO()
    var recipe Recipe

	if err := c.BindJSON(&recipe); err != nil {
        c.Status(400)
        return
    }
	update := bson.M{ "$set": recipe }
    if _, err := recipeDB.UpdateOne(ctx, bson.M{"_id": recipe.ID, "user_id": recipe.UserID}, update); err != nil {
	  fmt.Println(err)
	  c.Status(500)
	  return
	}
    c.Status(200)
    return
}

func recipeAddQueue(c *gin.Context) {
  session := sessions.Default(c)
	ctx := context.TODO()
	userIDString := session.Get("user")
	fmt.Println(userIDString)
	if userIDString == nil {
	  fmt.Println("session not found...")
	  c.Status(403)
	  return 
	}
	var recipe Recipe
	var device Device
	recipeID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := recipeDB.FindOne(ctx, bson.M{"_id": recipeID}).Decode(&recipe); err != nil {
	  fmt.Println(err)
		c.Status(404)
		return
	}
	if err := deviceDB.FindOne(ctx, bson.M{"user_id": userIDString}).Decode(&device); err != nil {
	  fmt.Println(err)
		c.Status(400)
		return
	}
	device.Recipe = recipe
	update := bson.M{ "$set": bson.M{ "recipe": recipe } }
	if _, err := deviceDB.UpdateOne(ctx, bson.M{"user_id": userIDString}, update); err != nil {
	  fmt.Println(err)
	  c.Status(500)
	  return
	}
	c.JSON(200, recipe)
}
