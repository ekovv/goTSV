package handler

import (
	"github.com/gin-gonic/gin"
	"goTSV/internal/service"
	"goTSV/internal/shema"
	"net/http"
)

type Handler struct {
	service service.Service
	engine  *gin.Engine
}

func NewHandler(service service.Service) *Handler {
	router := gin.Default()
	h := &Handler{
		service: service,
		engine:  router,
	}
	Route(router, h)
	return h
}

func (s *Handler) Start() {
	err := s.engine.Run("localhost:8080")
	if err != nil {
		return
	}
}

func (s *Handler) GetAll(c *gin.Context) {
	var req shema.Request
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
}
