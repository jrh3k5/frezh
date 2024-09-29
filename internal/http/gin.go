package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jrh3k5/frezh/internal/chatgpt"
	"github.com/jrh3k5/frezh/internal/http/handler"
	"github.com/jrh3k5/frezh/internal/http/handler/import/hellofresh"
	"github.com/jrh3k5/frezh/internal/http/handler/recipes"
	"github.com/jrh3k5/frezh/internal/ocr"
	recipesrepo "github.com/jrh3k5/frezh/internal/recipes"

	"github.com/jrh3k5/frezh/internal/http/handler/recipes/create"
)

func StartServer(chatpgptService chatgpt.Service, ocrProcessor ocr.Processor, recipesRepository recipesrepo.Repository) error {
	router := gin.Default()
	router.LoadHTMLGlob("internal/http/templates/*.tmpl")
	router.Static("/static", "internal/http/content/static")

	router.Use(gin.Recovery())

	router.GET("/", handler.HandleIndex)

	// HelloFresh
	router.GET("/import/hellofresh", hellofresh.HandleIndex)
	router.POST("/import/hellofresh", hellofresh.NewIngredientsUploadHandler(chatpgptService, ocrProcessor))

	// Recipes
	router.GET("/recipes/:id", recipes.NewRecipeGetHandler(recipesRepository))
	router.GET("/recipes/create", create.HandleIndex)
	router.POST("/recipes/create", create.NewRecipeCreationHandler(recipesRepository))

	err := router.Run(":8080")
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
