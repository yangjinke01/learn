package middleware

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

var store sessions.Store

func Session(engine *gin.Engine) {
	if store == nil {
		store, _ := redis.NewStore(10, "tcp", "10.0.50.16:6379", "", []byte("secret"))
		engine.Use(sessions.Sessions("mysession", store))
		engine.Use(func(c *gin.Context) {
			session := sessions.Default(c)
			session.Options(sessions.Options{MaxAge: 60, HttpOnly: true, Secure: false})
			//获取对应的token
			if session.Get("token") != "TOKEN" {
				fmt.Println(c.Request.RequestURI)
				if c.Request.RequestURI == "/user/login" {
					c.Next()
				} else {
					c.JSON(http.StatusBadRequest, "please login")
					c.Abort()
				}
			}
		})
	}
}
