package actions

import (
	"github.com/gin-gonic/gin"
	"microlog/input"
	"net/http"
	"strconv"
)

// All inputs
func Inputs(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"inputs.html",
		gin.H{
			"title":  "Inputs",
			"inputs": input.GetAllInputs(),
		},
	)
}

// Stop input
func StopInput(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	input := input.GetById(id)
	input.Shutdown()
	c.Redirect(http.StatusFound, "/inputs")
}

// Start input
func StartInput(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	input := input.GetById(id)
	input.Start()
	c.Redirect(http.StatusFound, "/inputs")
}
