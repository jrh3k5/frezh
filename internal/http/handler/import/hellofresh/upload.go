package hellofresh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jrh3k5/frezh/internal/chatgpt"
	frezhchatgpt "github.com/jrh3k5/frezh/internal/chatgpt"
	"github.com/jrh3k5/frezh/internal/http/handler/errors"
	"github.com/jrh3k5/frezh/internal/http/handler/recipes/create"
	"github.com/jrh3k5/frezh/internal/ocr"
)

func NewIngredientsUploadHandler(chatgptService frezhchatgpt.Service, ocrProcessor ocr.Processor) gin.HandlerFunc {
	return func(c *gin.Context) {
		ingredients, err := getIngredients(c, chatgptService, ocrProcessor)
		if err != nil {
			errors.HandleError(c, fmt.Errorf("failed to get ingredients: %w", err))

			return
		}

		steps, err := getRecipeSteps(c, chatgptService, ocrProcessor)
		if err != nil {
			errors.HandleError(c, fmt.Errorf("failed to get recipe steps: %w", err))

			return
		}

		recipeData := &create.RecipeData{
			Ingredients: ingredients,
			Steps:       steps,
		}

		serializedData, err := recipeData.Serialize()
		if err != nil {
			errors.HandleError(c, fmt.Errorf("failed to serialize recipe data: %w", err))

			return
		}

		c.Header("Location", "/recipes/create?recipe_base64="+base64.RawURLEncoding.EncodeToString(serializedData))
		c.Status(http.StatusSeeOther)
	}
}

func getIngredients(c *gin.Context, chatgptService chatgpt.Service, ocrProcessor ocr.Processor) ([]create.RecipeIngredient, error) {
	ingredientsFileHeader, err := c.FormFile("file_ingredients")
	if err != nil {
		return nil, fmt.Errorf("failed to get ingredients file from form data: %w", err)
	}

	file, err := ingredientsFileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	text, err := ocrProcessor.GetText(file)
	if err != nil {
		return nil, fmt.Errorf("failed to get text using OCR: %w", err)
	}

	slog.Info("OCR-processsed ingredients text", "text", text)

	answer, err := chatgptService.Ask(c, fmt.Sprintf("parse this three-column grid of ingredients where each cell is a picture of the ingredient, the quantity of the ingredient (for two and four people, respectively), and the name of the ingredient into a JSON array of objects where each object has three fields - ingredient_name, quantity_two_people, and quantity_four_people - with no Markdown formatting: %s", text))
	if err != nil {
		return nil, fmt.Errorf("failed to ask ChatGPT: %w", err)
	}

	slog.Info("ChatGPT-processsed ingredients answer", "answer", answer)

	var ingredients []chatgptIngredientList
	err = json.Unmarshal([]byte(answer), &ingredients)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	creationIngredients := make([]create.RecipeIngredient, 0, len(ingredients))
	for _, ingredient := range ingredients {
		var quantities []create.RecipeIngredientQuantity
		if ingredient.QuantityTwoPeople != "" {
			quantities = append(quantities, create.RecipeIngredientQuantity{
				ServingSize: 2,
				Value:       ingredient.QuantityTwoPeople,
			})
		}

		if ingredient.QuantityFourPeople != "" {
			quantities = append(quantities, create.RecipeIngredientQuantity{
				ServingSize: 4,
				Value:       ingredient.QuantityFourPeople,
			})
		}

		creationIngredients = append(creationIngredients, create.RecipeIngredient{
			Name:       ingredient.IngredientName,
			Quantities: quantities,
		})
	}

	return creationIngredients, nil
}

func getRecipeSteps(c *gin.Context, chatgptService chatgpt.Service, ocrProcessor ocr.Processor) ([]create.RecipeStep, error) {
	stepCountText := c.Request.FormValue("step_count")
	stepCount, err := strconv.ParseInt(stepCountText, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse step count: %w", err)
	}

	var steps []create.RecipeStep
	for i := 1; i <= int(stepCount); i++ {
		stepsFileHeader, err := c.FormFile("file_steps_" + strconv.Itoa(i))
		if err != nil {
			return nil, fmt.Errorf("failed to get steps file from form data at step %d: %w", i, err)
		}

		file, err := stepsFileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file for step %d: %w", i, err)
		}

		defer func() {
			_ = file.Close()
		}()

		text, err := ocrProcessor.GetText(file)
		if err != nil {
			return nil, fmt.Errorf("failed to get text using OCR: %w", err)
		}

		slog.Info("OCR-processsed steps text", "index", i, "text", text)

		answer, err := chatgptService.Ask(c, fmt.Sprintf("parse this set of instructions into a list of instructions with no Markdown formatting: %s", text))
		if err != nil {
			return nil, fmt.Errorf("failed to ask ChatGPT at step %d: %w", i, err)
		}

		slog.Info("ChatGPT-processsed ingredients answer", "index", i, "answer", answer)

		steps = append(steps, create.RecipeStep{
			StepText: answer,
		})
	}

	return steps, nil
}

type chatgptIngredientList struct {
	IngredientName     string `json:"ingredient_name"`
	QuantityTwoPeople  string `json:"quantity_two_people"`
	QuantityFourPeople string `json:"quantity_four_people"`
}
