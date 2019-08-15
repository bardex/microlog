package actions

import (
	"github.com/gin-gonic/gin"
	"microlog/listeners"
	"microlog/settings"
	"net/http"
	"strconv"
)

type inputJson struct {
	Id       int64  `json:id`
	Protocol string `json:protocol`
	Addr     string `json:addr`
	Active   bool   `json:active`
	Error    string `json:error`
}

// All inputs
func Inputs(c *gin.Context) {
	inputs, _ := settings.Inputs.GetAll()

	data := make([]inputJson, 0, len(inputs))

	for _, input := range inputs {
		data = append(data,
			inputJson{
				Id:       input.Id,
				Protocol: input.Protocol,
				Addr:     input.Addr,
				Active:   input.GetListener().IsActive(),
				Error:    input.GetListener().GetError(),
			})
	}

	c.JSON(
		http.StatusOK,
		data,
	)
}

// Stop input
func StopInput(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	input, err := settings.Inputs.GetOne(id)

	if err == nil {
		listener := input.GetListener()
		listener.Stop()
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
		listener := input.GetListener()
		listener.Start()
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
			"title":      "Add input",
			"extractors": listeners.Extractors,
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
