package config

import (
	"github.com/mathieu-gilloots/tcardgen/pkg/canvas/box"
	"github.com/mathieu-gilloots/tcardgen/pkg/canvas/fontfamily"
)

type DrawingConfig struct {
	Template string               `json:"template,omitempty"`
	Title    *MultiLineTextOption `json:"title,omitempty"`
	Category *TextOption          `json:"category,omitempty"`
	Info     *TextOption          `json:"info,omitempty"`
	Tags     *BoxTextsOption      `json:"tags,omitempty"`
}

type TextOption struct {
	Start      *Point           `json:"start,omitempty"`
	FgHexColor string           `json:"fgHexColor,omitempty"`
	FontSize   float64          `json:"fontSize,omitempty"`
	FontStyle  fontfamily.Style `json:"fontStyle,omitempty"`
	Separator  string           `json:"separator,omitempty"`
}

type MultiLineTextOption struct {
	TextOption
	MaxWidth    int  `json:"maxWidth,omitempty"`
	LineSpacing *int `json:"lineSpacing,omitempty"`
}

type BoxTextsOption struct {
	TextOption
	BgHexColor string    `json:"bgHexColor,omitempty"`
	BoxPadding *Padding  `json:"boxPadding,omitempty"`
	BoxSpacing *int      `json:"boxSpacing,omitempty"`
	BoxAlign   box.Align `json:"boxAlign,omitempty"`
}

type Point struct {
	X int `json:"px"`
	Y int `json:"py"`
}

type Padding struct {
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
}
