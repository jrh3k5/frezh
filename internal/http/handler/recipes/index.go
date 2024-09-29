package recipes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jrh3k5/frezh/internal/http/handler/errors"
	"github.com/jrh3k5/frezh/internal/recipes"
)

func NewRecipeGetHandler(recipesRepository recipes.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		recipeID := c.Param("id")
		recipe, err := recipesRepository.GetRecipe(c.Request.Context(), recipeID)
		if err != nil {
			errors.HandleError(c, fmt.Errorf("failed to get recipe: %w", err))

			return
		} else if recipe == nil {
			c.Status(http.StatusNotFound)

			return
		}

		ingredientNames := make([]string, 0, len(recipe.Ingredients))
		ingredientQuantities := make([][]string, 0, len(recipe.Ingredients))
		for i := 0; i < len(recipe.Ingredients); i++ {
			ingredient := recipe.Ingredients[i]

			ingredientNames = append(ingredientNames, ingredient.Name)

			quantities := make([]string, 0, len(ingredient.Quantities))
			for j := 0; j < len(ingredient.Quantities); j++ {
				quantity := ingredient.Quantities[j]

				quantities = append(quantities, quantity.Quantity)
			}

			ingredientQuantities = append(ingredientQuantities, quantities)
		}

		stepsText := make([]string, 0, len(recipe.Steps))
		for _, step := range recipe.Steps {
			stepsText = append(stepsText, step.StepText)
		}

		servingSizes := recipe.GetDistinctServingSizes()
		servingSizeLabels := make([]string, 0, len(servingSizes))
		for _, servingSize := range servingSizes {
			servingSizeLabels = append(servingSizeLabels, fmt.Sprintf("%d", servingSize))
		}

		c.HTML(http.StatusOK, "recipes_index_individual.tmpl", gin.H{
			"recipe_name":           recipe.Name,
			"ingredient_names":      ingredientNames,
			"ingredient_quantities": ingredientQuantities,
			"serving_sizes":         servingSizes,
			"serving_size_labels":   servingSizeLabels,
			"steps":                 stepsText,
		})
	}
}
