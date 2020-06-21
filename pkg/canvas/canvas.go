package canvas

import (
	"bytes"
	"image"
	"image/draw"

	"github.com/google/martian/log"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/Ladicle/tcardgen/pkg/canvas/fontfamily"
)

func CreateCanvasFromImage(filename string) (*Canvas, error) {
	tpl, err := LoadFromFile(filename)
	if err != nil {
		return nil, err
	}
	// draw background image
	dst := image.NewRGBA(tpl.Bounds())
	draw.Draw(dst, dst.Bounds(), tpl, image.Point{}, draw.Over)

	return &Canvas{
		dst: dst,
		fdr: &font.Drawer{Dst: dst, Src: image.Black, Dot: fixed.Point26_6{}},
	}, nil
}

type Canvas struct {
	dst *image.RGBA
	fdr *font.Drawer

	maxWidth  int
	lineSpace int
}

// SaveAsPNG saves this canvas as a PNG file into the specified path.
func (c *Canvas) SaveAsPNG(filename string) error {
	return SaveAsPNG(filename, c.dst)
}

// DrawTextAtPoint draws text on this canvas at the specified point.
func (c *Canvas) DrawTextAtPoint(text string, x, y int, opts ...textDrawOption) error {
	for _, f := range opts {
		if err := f(c); err != nil {
			return err
		}
	}

	if c.fdr.Face == nil {
		log.Errorf("Face is nil: %+v", c.fdr)
	}

	// dot.y points baseline of text
	c.fdr.Dot.Y = fixed.I(y) + c.fdr.Face.Metrics().Height
	c.fdr.Dot.X = fixed.I(x)

	if c.maxWidth == 0 {
		c.fdr.DrawString(text)
		return nil
	}

	c.drawMultiLineText(text)
	return nil
}

func (c *Canvas) drawMultiLineText(text string) {
	var (
		x      = c.fdr.Dot.X
		rtext  = []rune(text)
		length = len(rtext)

		lbuf bytes.Buffer
		wbuf bytes.Buffer
	)
	for i := 0; i < length; i++ {
		r := rtext[i]
		switch {
		case spaceChar(r):
			// noop
		case oneByteChar(r) || startBracket(r):
			wbuf.WriteRune(r)
			if (i + 1) < length {
				continue
			}
		default:
			wbuf.WriteRune(r)
			if (i+1) < length && endChar(rtext[i+1]) {
				wbuf.WriteRune(rtext[i+1])
				i++
			}
		}
		lbuf.Write(wbuf.Bytes())

		adv := c.fdr.MeasureBytes(lbuf.Bytes())
		if adv <= fixed.I(c.maxWidth) {
			wbuf.Reset()
			if (i + 1) < length {
				continue
			}
		}

		c.fdr.DrawBytes(lbuf.Bytes()[:lbuf.Len()-wbuf.Len()])
		c.fdr.Dot.X = x
		c.fdr.Dot.Y += c.fdr.Face.Metrics().Height + fixed.I(c.lineSpace)

		lbuf.Reset()
		lbuf.Write(wbuf.Bytes())
		wbuf.Reset()
	}
}

type textDrawOption func(*Canvas) error

// FontFace sets font face.
func FontFace(ff font.Face) textDrawOption {
	return func(c *Canvas) error {
		c.fdr.Face = ff
		return nil
	}
}

// FontFaceFromFFA sets font face from FontFamily.
func FontFaceFromFFA(ffa *fontfamily.FontFamily, style fontfamily.Style, size float64) textDrawOption {
	return func(c *Canvas) error {
		ff, err := ffa.NewFace(style, size)
		if err != nil {
			return err
		}
		c.fdr.Face = ff
		return nil
	}
}

// FgColor sets foreground color.
func FgColor(color *image.Uniform) textDrawOption {
	return func(c *Canvas) error {
		c.fdr.Src = color
		return nil
	}
}

// FgHexColor sets foreground color hex.
func FgHexColor(hex string) textDrawOption {
	return func(c *Canvas) error {
		color, err := Hex(hex)
		if err != nil {
			return err
		}
		c.fdr.Src = color
		return nil
	}
}

// MaxWidth sets maximum width of text.
// If the full text width exceeds the limit, drawer adds line breaks.
func MaxWidth(max int) textDrawOption {
	return func(c *Canvas) error {
		c.maxWidth = max
		return nil
	}
}

// LineSpace sets line space (px) of multi-line text.
func LineSpace(px int) textDrawOption {
	return func(c *Canvas) error {
		c.lineSpace = px
		return nil
	}
}
