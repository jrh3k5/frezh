package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jrh3k5/frezh/internal/chatgpt"
	"github.com/jrh3k5/frezh/internal/http/handler"
	"github.com/jrh3k5/frezh/internal/http/handler/import/hellofresh"
	"github.com/jrh3k5/frezh/internal/ocr"
)

func StartServer(chatpgptService chatgpt.Service, ocrProcessor ocr.Processor) error {
	router := gin.Default()
	router.LoadHTMLGlob("internal/http/templates/*.tmpl")

	router.Use(gin.Recovery())

	router.GET("/", handler.HandleIndex)
	router.GET("/import/hellofresh", hellofresh.HandleIndex)
	router.POST("/import/hellofresh", hellofresh.NewIngredientsUploadHandler(chatpgptService, ocrProcessor))

	err := router.Run(":8080")
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
