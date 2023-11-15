package canvas

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

// LoadFromFile loads an image file and generate image.Image from it.
// Supported image types are JPEG and PNG.
func LoadFromFile(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", filename, err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

// SaveAsPNG saves image object as a PNG image.
func SaveAsPNG(filename string, img image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %q: %w", filename, err)
	}
	defer f.Close()
	return png.Encode(f, img)
}
