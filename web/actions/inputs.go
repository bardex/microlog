package actions

import (
	"github.com/gin-gonic/gin"
	"microlog/settings"
	"net/http"
	"strconv"
)

// All inputs
func Inputs(c *gin.Context) {
	inputs, _ := settings.Inputs.GetAll()

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
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	input, err := settings.Inputs.GetOne(id)

	if err == nil {
		input.GetListener().Stop()
		input.Enabled = 0
		settings.Inputs.Update(input)
	}

	c.Redirect(http.StatusFound, "/inputs")
}

// Start input
func StartInput(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	input, err := settings.Inputs.GetOne(id)

	if err == nil {
		input.GetListener().Start()
		input.Enabled = 1
		settings.Inputs.Update(input)
	}

	c.Redirect(http.StatusFound, "/inputs")
}

// Delete input
func DeleteInput(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	settings.Inputs.Delete(id)

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
	extractor := c.PostForm("extractor")
	addr := c.PostForm("address")

	newInput := settings.Input{
		Protocol:  protocol,
		Extractor: extractor,
		Addr:      addr,
	}

	repo := settings.Inputs

	repo.Add(&newInput)

	c.Redirect(http.StatusFound, "/inputs")
}
