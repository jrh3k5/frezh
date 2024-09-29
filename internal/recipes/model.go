package recipes

import "slices"

type Recipe struct {
	ID          string
	Name        string
	Ingredients []RecipeIngredient
	Steps       []RecipeStep
}

// GetDistinctServingSizes returns the list of distinct serving sizes in the recipe.
func (r *Recipe) GetDistinctServingSizes() []int {
	uniqueServingSizes := make(map[int]any)

	for _, ingredient := range r.Ingredients {
		for _, quantity := range ingredient.Quantities {
			uniqueServingSizes[quantity.ServingSize] = nil
		}
	}

	servingSizes := make([]int, 0, len(uniqueServingSizes))
	for servingSize := range uniqueServingSizes {
		servingSizes = append(servingSizes, servingSize)
	}

	slices.Sort(servingSizes)

	return servingSizes
}

type RecipeIngredient struct {
	Name       string
	Quantities []RecipeIngredientQuantity
}

type RecipeIngredientQuantity struct {
	ServingSize int
	Quantity    string
}

type RecipeStep struct {
	StepText string
}
