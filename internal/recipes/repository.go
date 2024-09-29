package recipes

import "context"

// Repository defines a means of storing and retrieving recipes.
type Repository interface {
	// GetRecipe retrieves a recipe for the given ID.
	// This returns nil if the recipe does not exist.
	GetRecipe(ctx context.Context, id string) (*Recipe, error)

	// GetRecipes retrieves all recipes.
	GetRecipes(ctx context.Context) ([]Recipe, error)

	// SaveRecipe persists a recipe.
	// This returns the ID of the recipe if saving is successful.
	SaveRecipe(ctx context.Context, request CreateRecipeRequest) (string, error)
}

type CreateRecipeRequest struct {
	Name        string
	Ingredients []RecipeIngredient
	Steps       []RecipeStep
}
