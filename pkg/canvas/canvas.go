package canvas

import (
	"bytes"
	"image"
	"image/draw"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/mathieu-gilloots/tcardgen/pkg/canvas/box"
	"github.com/mathieu-gilloots/tcardgen/pkg/canvas/fontfamily"
	"github.com/mathieu-gilloots/tcardgen/pkg/config"
)

func CreateCanvasFromImage(tpl image.Image) (*Canvas, error) {
	// draw background image
	dst := image.NewRGBA(tpl.Bounds())
	draw.Draw(dst, dst.Bounds(), tpl, image.Point{}, draw.Src)

	return &Canvas{
		dst: dst,
		fdr: &font.Drawer{Dst: dst, Src: image.Black, Dot: fixed.Point26_6{}},
	}, nil
}

type Canvas struct {
	dst *image.RGBA
	fdr *font.Drawer

	bgColor    *image.Uniform
	maxWidth   int
	lineSpace  int
	boxPadding config.Padding
	boxSpace   int
	boxAlign   box.Align
}

// SaveAsPNG saves this canvas as a PNG file into the specified path.
func (c *Canvas) SaveAsPNG(filename string) error {
	return SaveAsPNG(filename, c.dst)
}

// DrawTextAtPoint draws text on this canvas at the specified point.
func (c *Canvas) DrawTextAtPoint(text string, start config.Point, opts ...textDrawOption) error {
	for _, f := range opts {
		if err := f(c); err != nil {
			return err
		}
	}

	// dot.y points baseline of text
	c.fdr.Dot.Y = fixed.I(start.Y) + c.fdr.Face.Metrics().Height
	c.fdr.Dot.X = fixed.I(start.X)

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

		wbuf.WriteRune(r)

		switch {
		case spaceChar(r):
			// noop
		case oneByteChar(r) || startBracket(r):
			if (i + 1) < length {
				continue
			}
		case (i+1) < length && endChar(rtext[i+1]):
			wbuf.WriteRune(rtext[i+1])
			i++
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

	if len(lbuf.Bytes()) != 0 {
		c.fdr.DrawBytes(lbuf.Bytes()[:lbuf.Len()-wbuf.Len()])
	}
}

func (c *Canvas) DrawBoxTexts(texts []string, start config.Point, opts ...textDrawOption) error {
	for _, f := range opts {
		if err := f(c); err != nil {
			return err
		}
	}

	p := image.Pt(start.X, start.Y)
	if c.boxAlign == box.AlignRight {
		n := len(texts)
		p.X -= c.boxPadding.Left*n + c.boxPadding.Right*n + c.boxSpace*(n-1) +
			c.fdr.MeasureString(strings.Join(texts, "")).Round()
	}

	fm := c.fdr.Face.Metrics()
	fh := fm.Height
	rect := image.Rect(0, start.Y, 0, start.Y+fh.Round()+c.boxPadding.Top+c.boxPadding.Bottom+fm.Descent.Round())

	for _, s := range texts {
		fw := c.fdr.MeasureString(s)
		rect.Min.X = p.X
		rect.Max.X = p.X + fw.Round() + c.boxPadding.Left + c.boxPadding.Right
		draw.Draw(c.dst, rect, c.bgColor, p, draw.Src)

		c.fdr.Dot.X = fixed.I(p.X + c.boxPadding.Left)
		c.fdr.Dot.Y = fixed.I(p.Y+c.boxPadding.Top-1) + fh
		c.fdr.DrawString(s)

		p.X = rect.Max.X + c.boxSpace
	}
	return nil
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

// BgColor sets background color.
func BgColor(color *image.Uniform) textDrawOption {
	return func(c *Canvas) error {
		c.bgColor = color
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

// BgHexColor sets background color hex.
func BgHexColor(hex string) textDrawOption {
	return func(c *Canvas) error {
		color, err := Hex(hex)
		if err != nil {
			return err
		}
		c.bgColor = color
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

// LineSpace sets line space(px) of multi-line text.
func LineSpacing(px int) textDrawOption {
	return func(c *Canvas) error {
		c.lineSpace = px
		return nil
	}
}

// BoxPadding sets box padding(px).
func BoxPadding(bp config.Padding) textDrawOption {
	return func(c *Canvas) error {
		c.boxPadding = bp
		return nil
	}
}

// BoxSpacing sets box spacing(px).
func BoxSpacing(px int) textDrawOption {
	return func(c *Canvas) error {
		c.boxSpace = px
		return nil
	}
}

// BoxAlign sets box align.
func BoxAlign(align box.Align) textDrawOption {
	return func(c *Canvas) error {
		c.boxAlign = align
		return nil
	}
}
