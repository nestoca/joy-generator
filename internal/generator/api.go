package generator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Output struct {
	Parameters []Result `json:"parameters"`
}

type GetParamsResponse struct {
	Output Output `json:"output"`
}

type GetParamsRequest struct {
	ApplicationSetName string `json:"applicationSetName"`
	Input              struct {
		Parameters map[string]string `json:"parameters"`
	} `json:"input"`
}

type API struct {
	Logger    zerolog.Logger
	Generator *Generator
}

func (api API) HandleGetParams(c *gin.Context) {
	var request GetParamsRequest
	if err := c.BindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
		return
	}

	results, err := api.Generator.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "failed to generate results",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetParamsResponse{Output: Output{Parameters: results}})
}
