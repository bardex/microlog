package actions

import (
	"github.com/gin-gonic/gin"
	"microlog/input"
	"net/http"
)

type inputViewModel struct {
	Protocol string
	Address  string
	Active   bool
	Disabled bool
	HasError bool
	Error    string
	ShowUrl  string
}

func Inputs(c *gin.Context) {

	var inputs = make([]inputViewModel, 0, len(input.GetAllInputs()))

	for _, row := range input.GetAllInputs() {
		inputs = append(inputs, inputViewModel{
			Protocol: row.GetProtocol(),
			Address:  row.GetAddr(),
			Active:   row.IsActive(),
			Disabled: !row.IsEnabled(),
			HasError: row.HasError(),
			Error:    row.GetError(),
			ShowUrl:  "/search",
		})
	}

	c.HTML(
		http.StatusOK,
		"inputs.html",
		gin.H{
			"title":  "Inputs",
			"inputs": inputs,
		},
	)
}
