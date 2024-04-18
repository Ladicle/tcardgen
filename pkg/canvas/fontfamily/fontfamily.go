package fontfamily

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type Style string

const (
	Thin    = "Thin"
	Light   = "Light"
	Regular = "Regular"
	Medium  = "Medium"
	Bold    = "Bold"
	Black   = "Black"
)

const (
	TrueTypeFontExt = ".ttf"
)

// LoadFromDir loads files and return FontFamily object from the specified directory.
// The directory name is used as a family name, and all font files in it are identified as part
// of the same font family.  Each filename must follows this `<name>-<style>.ttf`naming rule.
func LoadFromDir(dir string) (*FontFamily, error) {
	finfos, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	fs := NewFontFamily(filepath.Base(dir))
	for _, finfo := range finfos {
		fn := finfo.Name()
		ext := filepath.Ext(fn)
		if ext != TrueTypeFontExt {
			// skip non TTF file
			continue
		}

		name := fn[:len(fn)-len(ext)]
		ss := strings.Split(name, "-")
		if len(ss) != 2 {
			return nil, fmt.Errorf("failed to parse %q name", fn)
		}

		if err := fs.LoadFont(filepath.Join(dir, fn), Style(ss[1])); err != nil {
			return nil, err
		}
	}
	return fs, nil
}

// NewFontFamily initialize a FontFamily object and return it.
func NewFontFamily(name string) *FontFamily {
	return &FontFamily{
		Name:  name,
		fonts: make(map[Style]*truetype.Font),
	}
}

type FontFamily struct {
	Name  string
	fonts map[Style]*truetype.Font
}

// LoadFont loads TrueType font from a file.
func (fs *FontFamily) LoadFont(filename string, style Style) error {
	if filepath.Ext(filename) != TrueTypeFontExt {
		return fmt.Errorf("%q is not TrueTypeFont format", filepath.Base(filename))
	}
	fb, err := os.ReadFile(filename)
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

// NewFace creates a new font face with size option.
func (fs *FontFamily) NewFace(style Style, size float64) (font.Face, error) {
	f, ok := fs.fonts[style]
	if !ok {
		return nil, fmt.Errorf("this font family does not contain %q style font", style)
	}
	return truetype.NewFace(f, &truetype.Options{Size: size}), nil
}
