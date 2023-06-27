package user

import (
	"example/web-service-gin/controllers"
	"example/web-service-gin/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

var db = middleware.DBConnection()

type User struct {
	Name string
	Age  int
}

type DOModelUser struct {
	gorm.Model
	User
}

type Controller struct {
	controllers.BaseController
}

func (contr *Controller) Login(c *gin.Context) {
	session := sessions.Default(c)
	//value为服务端创建的token
	session.Set("token", "TOKEN")
	session.Save()
	c.JSON(http.StatusOK, gin.H{"login": "success"})
}

func (contr *Controller) List(c *gin.Context) {
	var user DOModelUser
	db.First(&user, 1)
	c.JSON(http.StatusOK, user)
}

func (contr *Controller) Add(ctx *gin.Context) {
	db.AutoMigrate(&DOModelUser{})
	user := DOModelUser{User: User{Name: "Ting", Age: 26}}
	db.Create(&user)
	ctx.JSON(http.StatusOK, user)
}
