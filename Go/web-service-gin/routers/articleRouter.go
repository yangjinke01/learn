package routers

import (
	"example/web-service-gin/controllers/article"
	"github.com/gin-gonic/gin"
)

func ArticleRouterRegister(r *gin.Engine) {
	g := r.Group("/article")
	{
		controller:=article.Controller{}
		g.GET("/", controller.List)
		g.POST("/add", controller.Add)
	}
}
