package ingredients

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jrh3k5/frezh/http/handler/errors"
	"github.com/jrh3k5/frezh/ocr"
)

func HandleIngredientsUpload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		errors.HandleError(c, fmt.Errorf("failed to get file: %w", err))

		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		errors.HandleError(c, fmt.Errorf("failed to open file: %w", err))

		return
	}

	text, err := (&ocr.Gosseract{}).GetText(file)
	if err != nil {
		errors.HandleError(c, fmt.Errorf("failed to get text: %w", err))

		return
	}

	// TODO: return as an actual webpage
	c.JSON(http.StatusOK, gin.H{"text": text})
}
