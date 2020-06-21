package image

import (
	"fmt"
	"image"
	"image/color"
)

// Hex create image.Uniform from the specified color hex.
func Hex(hex string) (*image.Uniform, error) {
	var r, g, b uint8
	n, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return nil, err
	}
	if n != 3 {
		return nil, fmt.Errorf("failed to parse %v as a hex color", hex)
	}
	return image.NewUniform(color.RGBA{r, g, b, 255}), nil
}
