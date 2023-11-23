package handler

import (
	"github.com/gin-gonic/gin"
	"goTSV/config"
	"goTSV/internal/domains"
	"goTSV/internal/shema"
	"net/http"
)

type Handler struct {
	service domains.Service
	engine  *gin.Engine
	config  config.Config
}

func NewHandler(service domains.Service, cnf config.Config) *Handler {
	router := gin.Default()
	h := &Handler{
		service: service,
		engine:  router,
		config:  cnf,
	}
	Route(router, h)
	return h
}

func (s *Handler) Start() {
	err := s.engine.Run(s.config.Host)
	if err != nil {
		return
	}
}

func (s *Handler) GetAll(c *gin.Context) {
	var r shema.Request
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	result, err := s.service.GetAll(r)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)

}
