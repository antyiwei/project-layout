package handler

import "github.com/gin-gonic/gin"

type Handler struct {
	// Service *service.Service
}

func (h *Handler) GetByID(c *gin.Context) {
	c.JSON(200, gin.H{"message": "user.GetByID: not implemented"})
}
