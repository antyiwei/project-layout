package router

import (
	"github.com/gin-gonic/gin"

	orderhandler "project-layout/internal/order/handler"
)

func registerOrderRoutes(v1 *gin.RouterGroup, h *orderhandler.Handler) {
	if h == nil {
		return
	}
	orders := v1.Group("/orders")
	orders.POST("", h.Create)
}
