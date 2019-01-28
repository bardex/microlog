package web

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Start() {
	router = gin.Default()
	router.LoadHTMLGlob("web/templates/*")
	initRoutes()
	router.Run()
}