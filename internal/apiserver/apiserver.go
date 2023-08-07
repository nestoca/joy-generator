package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/nestoca/joy-generator/internal/generator"
)

type GeneratorOutput struct {
	Parameters []*generator.Result `json:"parameters"`
}

type GetParamsResponse struct {
	Output GeneratorOutput `json:"output"`
}

type GeneratorInput struct {
	Parameters map[string]string `json:"parameters"`
}

type GetParamsRequest struct {
	ApplicationSetName string         `json:"applicationSetName"`
	Input              GeneratorInput `json:"input"`
}

type ApiServer struct {
	generator *generator.Generator
}

func New(g *generator.Generator) *ApiServer {
	return &ApiServer{
		generator: g,
	}
}

func (s *ApiServer) Run() error {
	r := gin.Default()

	r.GET("/api/v1/health", s.Health)
	//goland:noinspection SpellCheckingInspection
	r.GET("/api/v1/getparams.execute", s.GetParamsExecute)

	return r.Run()
}

func (s *ApiServer) Health(c *gin.Context) {
	if err := s.generator.Status(); err != nil {
		log.Error().Err(err).Msg("health check failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *ApiServer) GetParamsExecute(c *gin.Context) {
	body := &GetParamsRequest{}
	err := c.BindJSON(body)
	if err != nil {
		log.Debug().Err(err).Msg("invalid request body received")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"error":  "invalid request body",
			"detail": err.Error(),
		})
		return
	}

	results, err := s.generator.Run()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate results")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":   http.StatusInternalServerError,
			"error":  "failed to generate results",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, GetParamsResponse{
		Output: GeneratorOutput{
			Parameters: results,
		},
	})
}
