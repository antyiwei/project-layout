package router

import (
	"github.com/gin-gonic/gin"

	orderhandler "project-layout/internal/order/handler"
	userhandler "project-layout/internal/user/handler"
)

type Handlers struct {
	User  *userhandler.Handler
	Order *orderhandler.Handler
}

func Engine(h Handlers) *gin.Engine {
	r := gin.Default()
	Register(r, h)
	return r
}

func Register(r *gin.Engine, h Handlers) {
	v1 := r.Group("/api/v1")
	registerUserRoutes(v1, h.User)
	registerOrderRoutes(v1, h.Order)
}
