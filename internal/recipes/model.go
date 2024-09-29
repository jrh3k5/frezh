package recipes

type Recipe struct {
	ID          string
	Name        string
	Ingredients []RecipeIngredient
	Steps       []RecipeStep
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
