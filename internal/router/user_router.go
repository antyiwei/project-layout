package router

import (
	"github.com/gin-gonic/gin"

	userhandler "project-layout/internal/user/handler"
)

func registerUserRoutes(v1 *gin.RouterGroup, h *userhandler.Handler) {
	if h == nil {
		return
	}
	users := v1.Group("/users")
	users.GET("/:id", h.GetByID)
}
