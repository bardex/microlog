package actions

import (
	"github.com/gin-gonic/gin"
	"microlog/input"
	"microlog/settings"
	"net/http"
	"strconv"
)

// All inputs
func Inputs(c *gin.Context) {
	repo := settings.Inputs{}
	inputs, _ := repo.GetAll()
	c.HTML(
		http.StatusOK,
		"inputs.html",
		gin.H{
			"title":  "Inputs",
			"inputs": inputs,
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

// Add input
func AddInput(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"add_input.html",
		gin.H{
			"title": "Add input",
		},
	)
}

// Create new input
func CreateInput(c *gin.Context) {

	protocol := c.PostForm("protocol")
	addr := c.PostForm("address")

	newInput := settings.Input{
		Protocol: protocol,
		Addr:     addr,
	}

	repo := settings.Inputs{}

	repo.Add(&newInput)

	c.Redirect(http.StatusFound, "/inputs")
}
