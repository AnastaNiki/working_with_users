package handler

import (
	"github.com/gin-gonic/gin"
	"working_with_users/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	user := router.Group("/user")
	{
		user.POST("/create", h.createUser)
		user.POST("/update", h.updateUser)
		user.POST("/archive", h.archiveUser)
		user.POST("/find", h.findUsers)

		user.GET("/:id", h.getUserById)
	}
	return router
}
