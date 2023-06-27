package article

import (
	"net/http"

	"example/web-service-gin/controllers"
	"github.com/gin-gonic/gin"
)

type Article struct {
	Name   string
	Author string
}

type Controller struct {
	controllers.BaseController
}

func (contr *Controller) List(c *gin.Context) {
	articles := []Article{
		{"The secret garden", "Jack"},
		{"The little prince", "Ting"},
	}
	c.JSON(http.StatusOK, articles)
}

func (contr *Controller) Add(ctx *gin.Context) {
	article := Article{Name: "The secret garden", Author: "Jack"}
	ctx.JSON(http.StatusOK, article)
}
