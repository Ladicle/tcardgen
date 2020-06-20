package image

import (
	"golang.org/x/image/font"
)

func DrawText(dr *font.Drawer, text string, maxWidth int) error {
	dr.DrawString(text)
	return nil
}
