package actions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Search(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"search.html",
		gin.H{
			"title": "Search",
		},
	)
}
