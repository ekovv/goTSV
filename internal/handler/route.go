package handler

import "github.com/gin-gonic/gin"

func Route(c *gin.Engine, h *Handler) {
	c.POST("/api/all", h.GetAll)
}
