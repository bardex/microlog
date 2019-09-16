package web

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Start() {
	//gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	router.Delims("[[", "]]")
	router.LoadHTMLGlob("web/templates/*")
	initRoutes()
	router.Run()
}
