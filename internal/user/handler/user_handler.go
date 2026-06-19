package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"project-layout/internal/user/service"
)

type Handler struct {
	Service *service.Service
}

func (h *Handler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user, err := h.Service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
