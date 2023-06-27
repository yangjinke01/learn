package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct{}

func (bc *BaseController) Success(c *gin.Context) {
	c.String(http.StatusOK, "success")
}

func (bc *BaseController) Error(c *gin.Context) {
	c.String(http.StatusBadRequest, "error")
}
