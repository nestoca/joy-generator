package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/nestoca/joy-generator/internal/generator"
	"github.com/nestoca/joy-generator/internal/github"
)

type HandlerParams struct {
	pluginToken string
	logger      zerolog.Logger
	repo        *github.Repo
	generator   *generator.Generator
}

func Handler(params HandlerParams) http.Handler {
	engine := gin.New()

	engine.Use(func(c *gin.Context) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			params.logger.Err(fmt.Errorf("%v", err)).Msg("recovered from panic")

			if c.Writer.Written() {
				return
			}

			c.JSON(500, gin.H{"error": err})
		}()
	})

	engine.Use(func(c *gin.Context) {
		start := time.Now()

		recorder := ErrorRecorder{
			ResponseWriter: c.Writer,
			buffer:         bytes.Buffer{},
		}

		c.Writer = &recorder

		c.Next()

		event := func() *zerolog.Event {
			if err := recorder.buffer.String(); err != "" {
				return params.logger.Err(errors.New(err))
			}
			return params.logger.Info()
		}()

		event.
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("code", c.Writer.Status()).
			Dur("elapsed", time.Since(start)).
			Msg("served request")
	})

	engine.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.GET("/api/v1/readiness", func(c *gin.Context) {
		if err := params.repo.Pull(); err != nil {
			c.JSON(500, gin.H{
				"status": "error",
				"detail": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{"status": "ok"})
	})

	generatorAPI := generator.API{
		Logger:    params.logger,
		Generator: params.generator,
	}

	engine.GET(
		"/api/v1/getparams.execute",
		func(c *gin.Context) {
			if c.GetHeader("Authorization") != "Bearer "+params.pluginToken {
				c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			}
		},
		generatorAPI.HandleGetParams,
	)

	return engine.Handler()
}

type ErrorRecorder struct {
	gin.ResponseWriter
	buffer bytes.Buffer
}

func (recorder *ErrorRecorder) Write(data []byte) (int, error) {
	if recorder.Status() >= 400 {
		_, _ = recorder.buffer.Write(data)
	}
	return recorder.ResponseWriter.Write(data)
}
