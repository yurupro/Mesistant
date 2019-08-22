package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`            // user id
	Name     string             `bson:"name" json:"name"`         // ユーザー名
	Mail     string             `bson:"mail" json:"mail"`         // ユーザーメルアド
	Password string             `bson:"password" json:"password"` // ユーザーメルアド
	Sex      string             `bson:"sex" json:"sex"`           // 性別
}

func userLogout(c *gin.Context) {
  session := sessions.Default(c)
  session.Clear()
	if err := session.Save(); err != nil {
		c.Status(400)
		return
	}
	c.Status(200)
}

func userLogin(c *gin.Context) {
  session := sessions.Default(c)
	user := new(User)
	user.ID = primitive.NilObjectID
	userSessionID := session.Get("user")
	fmt.Println("userSessionID: ")
	fmt.Println(userSessionID)
	if userSessionID != nil {
		c.Data(400, "text/plain", []byte("You already logined."))
		fmt.Println("You already logined")
		return
	}

	ctx := context.Background()
	type LoginReqJSON struct {
		Mail     string `json:"mail"`
		Password string `json:"password"`
	}
	var login LoginReqJSON
	if err := c.BindJSON(&login); err != nil {
		fmt.Println(err)
		c.Status(400)
		return
	}
	fmt.Println("- Req")
	fmt.Println(login)
	filter := bson.M{"mail": login.Mail, "password": login.Password}
	fmt.Println(filter)
	if err := userDB.FindOne(ctx, filter).Decode(&user); err != nil {
		c.Status(400)
		fmt.Println(err)
		return
	}

	fmt.Println("- Found user")
	fmt.Println(user)
	fmt.Println(user.ID)
	session.Set("user", user.ID.Hex())
	session.Save()
	c.JSON(200, user)
}

func userAdd(c *gin.Context) {
	var user User
	user.ID = primitive.NewObjectID()
	c.BindJSON(&user)
	fmt.Println(user)
	if _, err := userDB.InsertOne(context.TODO(), user); err != nil {
		c.Status(500)
	}
	c.Status(200)
}

func userGet(c *gin.Context) {
	ctx := context.Background()
	var user User
	filter := bson.M{"_id": c.Param("id")}
	if err := userDB.FindOne(ctx, filter).Decode(&user); err != nil {
		c.Status(500)
	}
	c.JSON(200, user)
}
