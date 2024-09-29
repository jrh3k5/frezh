package create

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang/snappy"
	"github.com/jrh3k5/frezh/internal/http/handler/errors"
	"gopkg.in/yaml.v2"
)

func HandleIndex(c *gin.Context) {
	recipeDataBase64 := c.Request.URL.Query().Get("recipe_base64")

	var recipeData *RecipeData
	if recipeDataBase64 == "" {
		recipeData = &RecipeData{}
	} else {
		recipeBytes, err := base64.RawURLEncoding.DecodeString(recipeDataBase64)
		if err != nil {
			errors.HandleError(c, fmt.Errorf("failed to decode recipe base64: %w", err))

			return
		}

		recipeData, err = DeserializeRecipeData(recipeBytes)
		if err != nil {
			errors.HandleError(c, fmt.Errorf("failed to deserialize recipe: %w", err))

			return
		}
	}

	servingSizes := recipeData.GetDistinctServingSizes()
	ingredientNames := make([]string, 0, len(recipeData.Ingredients))
	ingredientQuantities := make([][]string, 0, len(recipeData.Ingredients))
	for _, ingredient := range recipeData.Ingredients {
		ingredientNames = append(ingredientNames, ingredient.Name)

		quantities := make([]string, 0, len(servingSizes))
		for i := 0; i < len(servingSizes); i++ {
			quantities = append(quantities, ingredient.GetValueForServingSize(servingSizes[i]))
		}

		ingredientQuantities = append(ingredientQuantities, quantities)
	}

	servingSizeLabels := make([]string, 0, len(servingSizes))
	for i := 0; i < len(servingSizes); i++ {
		servingSizeLabels = append(servingSizeLabels, fmt.Sprintf("Quantity (%d servings)", servingSizes[i]))
	}

	stepsText := make([]string, 0, len(recipeData.Steps))
	for _, step := range recipeData.Steps {
		stepsText = append(stepsText, step.StepText)
	}

	c.HTML(http.StatusOK, "recipes_create_index.tmpl", gin.H{
		"ingredient_names":      ingredientNames,
		"ingredient_quantities": ingredientQuantities,
		"serving_size_labels":   servingSizeLabels,
		"steps":                 stepsText,
	})
}

type RecipeData struct {
	Name        string             `json:"name"`
	Ingredients []RecipeIngredient `json:"ingredients"`
	Steps       []RecipeStep       `json:"steps"`
}

// GetDistinctServingSizes returns the list of distinct serving sizes in the recipe.
func (r *RecipeData) GetDistinctServingSizes() []int {
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

// Deserialize deserializes the recipe data contained within the given byte array.
func DeserializeRecipeData(data []byte) (*RecipeData, error) {
	yamlBytes, err := snappy.Decode(nil, data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal recipe: %w", err)
	}

	var recipeData RecipeData
	err = yaml.Unmarshal(yamlBytes, &recipeData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal recipe: %w", err)
	}

	return &recipeData, nil
}

// Serialize serializes the recipe data contained within this object.
func (r *RecipeData) Serialize() ([]byte, error) {
	yamlBytes, err := yaml.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal recipe: %w", err)
	}

	return snappy.Encode(nil, yamlBytes), nil
}

type RecipeIngredient struct {
	Name       string                     `json:"name"`
	Quantities []RecipeIngredientQuantity `json:"quantities"`
}

// GetValueForServingSize returns the value of the ingredient for the given serving size.
// If no quantity is found for the given serving size, an empty string is returned.
func (r *RecipeIngredient) GetValueForServingSize(servingSize int) string {
	for _, quantity := range r.Quantities {
		if quantity.ServingSize == servingSize {
			return quantity.Value
		}
	}

	return ""
}

type RecipeIngredientQuantity struct {
	ServingSize int    `json:"serving_size"`
	Value       string `json:"value"`
}

type RecipeStep struct {
	StepText string `json:"step_text"`
}
