package recipes

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type InMemoryRepository struct {
	recipesMutex sync.RWMutex
	recipes      []Recipe
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		recipes: make([]Recipe, 0),
	}
}

func (r *InMemoryRepository) GetRecipe(_ context.Context, id string) (*Recipe, error) {
	r.recipesMutex.RLock()
	defer r.recipesMutex.RUnlock()

	for _, recipe := range r.recipes {
		if recipe.ID == id {
			return &recipe, nil
		}
	}

	return nil, nil
}

func (r *InMemoryRepository) GetRecipes(_ context.Context) ([]Recipe, error) {
	r.recipesMutex.RLock()
	defer r.recipesMutex.RUnlock()

	recipesCopy := make([]Recipe, len(r.recipes))
	copy(recipesCopy, r.recipes)

	return recipesCopy, nil
}

func (r *InMemoryRepository) SaveRecipe(_ context.Context, request CreateRecipeRequest) (string, error) {
	r.recipesMutex.Lock()
	defer r.recipesMutex.Unlock()

	recipeID := uuid.NewString()
	recipe := Recipe{
		ID:          recipeID,
		Name:        request.Name,
		Ingredients: request.Ingredients,
		Steps:       request.Steps,
	}

	r.recipes = append(r.recipes, recipe)
	return recipeID, nil
}
