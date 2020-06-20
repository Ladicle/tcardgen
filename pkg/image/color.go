package image

import (
	"fmt"
	"image/color"
)

func Hex(hex string) (color.RGBA, error) {
	var r, g, b uint8
	n, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return color.RGBA{}, err
	}
	if n != 3 {
		return color.RGBA{}, fmt.Errorf("failed to parse %v as a hex color", hex)
	}
	return color.RGBA{r, g, b, 255}, nil
}
