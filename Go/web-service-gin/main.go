package main

import (
	"example/web-service-gin/middleware"
	"example/web-service-gin/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	middleware.Session(engine)
	routers.UserRouterRegister(engine)
	routers.ArticleRouterRegister(engine)
	engine.Run(":8080")
}
