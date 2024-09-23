package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jrh3k5/frezh/http/handler"
	"github.com/jrh3k5/frezh/http/handler/ingredients"
)

func StartServer() error {
	router := gin.Default()
	router.LoadHTMLGlob("http/templates/*")

	router.Use(gin.Recovery())

	router.GET("/", handler.HandleIndex)
	router.POST("/ingredients/upload", ingredients.HandleIngredientsUpload)

	err := router.Run(":8080")
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
