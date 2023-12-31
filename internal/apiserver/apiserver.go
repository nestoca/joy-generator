package apiserver

import (
	"net/http"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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
	token     string
	generator *generator.Generator
}

func New(token string, g *generator.Generator) *ApiServer {
	return &ApiServer{
		token:     token,
		generator: g,
	}
}

func (s *ApiServer) Run() error {
	r := gin.New()

	r.Use(logger.SetLogger(logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
		return l.Output(gin.DefaultWriter).With().Logger()
	})))

	r.GET("/api/v1/health", s.Health)
	r.GET("/api/v1/readiness", s.Readiness)
	//goland:noinspection SpellCheckingInspection
	r.POST("/api/v1/getparams.execute", s.GetParamsExecute)

	return r.Run()
}

func (s *ApiServer) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *ApiServer) Readiness(c *gin.Context) {
	if err := s.generator.Status(); err != nil {
		log.Error().Err(err).Msg("readiness check failed")
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
	if c.GetHeader("Authorization") != "Bearer "+s.token {
		log.Debug().Msg("invalid token received")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":   http.StatusUnauthorized,
			"error":  "invalid token",
			"detail": "invalid token received",
		})
		return
	}

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
