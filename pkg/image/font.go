package image

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"
)

type FontStyle string

const (
	FontStyleThin    = "Thin"
	FontStyleLight   = "Light"
	FontStyleRegular = "Regular"
	FontStyleMedium  = "Medium"
	FontStyleBold    = "Bold"
	FontStyleBlack   = "Black"
)

const (
	TrueTypeFontExt = ".ttf"
)

func LoadFontFamilyFromDir(dir string) (*FontFamily, error) {
	finfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	fs := NewFontFamily()
	for _, finfo := range finfos {
		fn := finfo.Name()
		name := fn[:len(fn)-len(filepath.Ext(fn))]
		ss := strings.Split(name, "-")
		if len(ss) != 2 {
			return nil, fmt.Errorf("failed to parse %q name", fn)
		}
		if err := fs.LoadFont(filepath.Join(dir, fn), FontStyle(ss[1])); err != nil {
			return nil, err
		}
	}
	return fs, nil
}

func NewFontFamily() *FontFamily {
	return &FontFamily{fonts: make(map[FontStyle]*truetype.Font)}
}

type FontFamily struct {
	fonts map[FontStyle]*truetype.Font
}

// LoadFont loads font from a file
func (fs *FontFamily) LoadFont(filename string, style FontStyle) error {
	if filepath.Ext(filename) != TrueTypeFontExt {
		return fmt.Errorf("%q is not TrueTypeFont format", filepath.Base(filename))
	}
	fb, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	f, err := truetype.Parse(fb)
	if err != nil {
		return err
	}
	if f == nil {
		return errors.New("parsed font is nil")
	}
	fs.fonts[style] = f
	return nil
}

func (fs *FontFamily) GetFont(style FontStyle) *truetype.Font {
	return fs.fonts[style]
}
