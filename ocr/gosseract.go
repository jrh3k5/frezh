package ocr

import (
	"fmt"
	"io"

	"github.com/otiai10/gosseract/v2"
)

type Gosseract struct {
}

func (g *Gosseract) GetText(imageReader io.Reader) (string, error) {
	imageBytes, err := io.ReadAll(imageReader)
	if err != nil {
		return "", fmt.Errorf("failed to read image bytes: %w", err)
	}

	client := gosseract.NewClient()
	defer func() {
		_ = client.Close()
	}()

	client.SetImageFromBytes(imageBytes)

	imageText, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("failed to get image text: %w", err)
	}

	return imageText, nil
}
