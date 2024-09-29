package create

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jrh3k5/frezh/internal/http/handler/errors"
	"github.com/jrh3k5/frezh/internal/recipes"
)

func NewRecipeCreationHandler(recipesRepository recipes.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		ingredients, err := getIngredients(c)
		if err != nil {
			errors.HandleError(c, fmt.Errorf("failed to get ingredients: %w", err))

			return
		}

		steps, err := getSteps(c)
		if err != nil {
			errors.HandleError(c, fmt.Errorf("failed to get steps: %w", err))

			return
		}

		recipeID, err := recipesRepository.SaveRecipe(c.Request.Context(), recipes.CreateRecipeRequest{
			Name:        c.Request.FormValue("recipe_name"),
			Ingredients: ingredients,
			Steps:       steps,
		})

		if err != nil {
			errors.HandleError(c, fmt.Errorf("failed to save recipe: %w", err))

			return
		}

		c.Header("Location", "/recipes/"+recipeID)
		c.Status(http.StatusSeeOther)
	}
}

func getIngredients(c *gin.Context) ([]recipes.RecipeIngredient, error) {
	ingredientCount := 0
	if ingredientCountText := c.Request.FormValue("ingredient_count"); ingredientCountText != "" {
		if parsedCount, err := strconv.ParseInt(ingredientCountText, 10, 32); err != nil {
			return nil, fmt.Errorf("failed to parse ingredient count: %w", err)
		} else {
			ingredientCount = int(parsedCount)
		}
	}

	quantityCount := 0
	if quantityCountText := c.Request.FormValue("quantity_count"); quantityCountText != "" {
		if parsedCount, err := strconv.ParseInt(quantityCountText, 10, 32); err != nil {
			return nil, fmt.Errorf("failed to parse quantity count: %w", err)
		} else {
			quantityCount = int(parsedCount)
		}
	}

	recipeIngredients := make([]recipes.RecipeIngredient, 0, ingredientCount)
	for i := 0; i < ingredientCount; i++ {
		ingredientName := c.Request.FormValue(fmt.Sprintf("ingredients[%d].name", i))

		ingredientQuantities := make([]recipes.RecipeIngredientQuantity, 0, quantityCount)
		for j := 0; j < quantityCount; j++ {
			servingSizeText := c.Request.FormValue(fmt.Sprintf("ingredients[%d].quantities[%d].serving_size", i, j))
			servingSize, err := strconv.ParseInt(servingSizeText, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("failed to parse serving size '%s': %w", servingSizeText, err)
			}

			quantity := c.Request.FormValue(fmt.Sprintf("ingredients[%d].quantities[%d].quantity", i, j))

			ingredientQuantities = append(ingredientQuantities, recipes.RecipeIngredientQuantity{
				ServingSize: int(servingSize),
				Quantity:    quantity,
			})
		}

		recipeIngredients = append(recipeIngredients, recipes.RecipeIngredient{
			Name:       ingredientName,
			Quantities: ingredientQuantities,
		})
	}

	return recipeIngredients, nil
}

func getSteps(c *gin.Context) ([]recipes.RecipeStep, error) {
	stepCountText := c.Request.FormValue("step_count")
	stepCount, err := strconv.ParseInt(stepCountText, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse step count: %w", err)
	}

	var steps []recipes.RecipeStep
	for i := 0; i < int(stepCount); i++ {
		stepText := c.Request.FormValue(fmt.Sprintf("steps[%d].step_text", i))

		steps = append(steps, recipes.RecipeStep{
			StepText: stepText,
		})
	}

	return steps, nil
}
