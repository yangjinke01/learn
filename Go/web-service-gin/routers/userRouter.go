package routers

import (
	"example/web-service-gin/controllers/user"
	"github.com/gin-gonic/gin"
)

func UserRouterRegister(r *gin.Engine) {
	g := r.Group("/user")
	{
		controller := user.Controller{}
		g.GET("/", controller.List)
		g.GET("/add", controller.Add)
		g.GET("/login",controller.Login)
	}
}
