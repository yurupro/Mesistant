package main
import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sex int
const (
  male Sex = iota
  female
)

type User struct {
	ID          primitive.ObjectID  `bson:"_id" json:"id"`          // user id
	Name        string `bson:"name" json:"name"`        // ユーザー名
	Mail        string `bson:"mail" json:"mail"`        // ユーザーメルアド
	Password        string `bson:"password" json:"password"`        // ユーザーメルアド
	Sex string `bson:"sex" json:"sex"` // 性別
}

func userLogout(c *gin.Context) {
  session := sessions.Default(c)
  session.Clear()
  session.Save()
  c.Status(200)
}
  
func userLogin(c *gin.Context) {
  session := sessions.Default(c)
  user := new(User)
  user.ID = primitive.NilObjectID
  v := session.Get("user")
  if v != nil {
	c.Data(400, "text/plain", []byte("You already logined."))
	fmt.Println("You already logined")
	fmt.Println(v)
	return
  }
  
  ctx := context.Background()
  var login User
  c.BindJSON(&login)
  fmt.Println("- Req")
  fmt.Println(login)
  filter := bson.M{"mail": login.Mail, "password": login.Password}
  fmt.Println(filter)
  if err := userDB.FindOne(ctx, filter).Decode(user); err != nil{
	c.Status(400)
	fmt.Println(err)
	return
  }

  fmt.Println("- Found user")
  fmt.Println(user)
  session.Set("session", user.ID)
  session.Save()
  c.JSON(200, user)
}

func userAdd(c *gin.Context) {
  var user User
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
  if err := userDB.FindOne(ctx, filter).Decode(&user); err != nil{
	c.Status(500)
  }
  c.JSON(200, user)
}

