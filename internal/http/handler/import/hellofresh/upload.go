package hellofresh

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jrh3k5/frezh/internal/chatgpt"
	frezhchatgpt "github.com/jrh3k5/frezh/internal/chatgpt"
	"github.com/jrh3k5/frezh/internal/http/handler/errors"
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

		// TODO: return as an actual webpage
		c.JSON(http.StatusOK, gin.H{"ingredients": ingredients, "steps": steps})
	}
}

func getIngredients(c *gin.Context, chatgptService chatgpt.Service, ocrProcessor ocr.Processor) ([]chatgptIngredientList, error) {
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

	return ingredients, nil
}

func getRecipeSteps(c *gin.Context, chatgptService chatgpt.Service, ocrProcessor ocr.Processor) ([]string, error) {
	stepsFileHeader, err := c.FormFile("file_steps")
	if err != nil {
		return nil, fmt.Errorf("failed to get steps file from form data: %w", err)
	}

	file, err := stepsFileHeader.Open()
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

	slog.Info("OCR-processsed steps text", "text", text)

	answer, err := chatgptService.Ask(c, fmt.Sprintf("parse this three-column grid of 6 steps where each starts with a number identifying its order and it may have bulleted steps beneath it into a JSON array of strings, ignoring any block of text that comes before the first step: %s", text))
	if err != nil {
		return nil, fmt.Errorf("failed to ask ChatGPT: %w", err)
	}

	slog.Info("ChatGPT-processsed ingredients answer", "answer", answer)

	var steps []string
	err = json.Unmarshal([]byte(answer), &steps)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return steps, nil
}

type chatgptIngredientList struct {
	IngredientName     string `json:"ingredient_name"`
	QuantityTwoPeople  string `json:"quantity_two_people"`
	QuantityFourPeople string `json:"quantity_four_people"`
}
