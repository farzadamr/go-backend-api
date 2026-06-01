package router

import (
	"github.com/farzadamr/go-backend-api/api/handler"
	"github.com/gin-gonic/gin"
)

func Item(router *gin.RouterGroup) {
	h := handler.NewItemHandler()

	router.POST("/", h.Create)
	router.GET("/:id", h.GetByID)
	router.PATCH("/:id", h.Update)
	router.DELETE("/:id", h.Delete)
	router.GET("/", h.List)

	//router.PATCH("/:id/availability", h.SetAvailability)
}
