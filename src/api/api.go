package api

import (
	"fmt"

	"github.com/farzadamr/go-backend-api/api/router"
	"github.com/farzadamr/go-backend-api/config"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) error {
	r := gin.New()
	r.Use(gin.Logger())

	RegisterRoutes(r)

	err := r.Run(fmt.Sprintf(":%s", cfg.HTTP.Port))
	if err != nil {
		return err
	}
	return nil
}

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		items := v1.Group("/items")
		router.Item(items)
	}
}
