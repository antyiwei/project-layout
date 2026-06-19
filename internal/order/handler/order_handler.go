package handler

import "github.com/gin-gonic/gin"

type Handler struct {
	// Service *service.Service
}

func (h *Handler) Create(c *gin.Context) {
	c.JSON(200, gin.H{"message": "order.Create: not implemented"})
}
