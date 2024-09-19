package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeapotController struct {
}

func (p TeapotController) Teapot(c *gin.Context) {
	c.JSON(http.StatusTeapot, nil)
}
