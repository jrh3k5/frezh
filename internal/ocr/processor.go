package ocr

import "io"

type Processor interface {
	GetText(imageReader io.Reader) (string, error)
}
