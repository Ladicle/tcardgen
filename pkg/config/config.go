package config

import (
	"github.com/Ladicle/tcardgen/pkg/canvas/box"
	"github.com/Ladicle/tcardgen/pkg/canvas/fontfamily"
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
	TimeFormat string           `json:"timeFormat,omitempty"`
	Enabled    *bool            `json:"enabled,omitempty"`
}

type MultiLineTextOption struct {
	TextOption
	MaxWidth    int   `json:"maxWidth,omitempty"`
	LineSpacing *int  `json:"lineSpacing,omitempty"`
	Enabled     *bool `json:"enabled,omitempty"`
}

type BoxTextsOption struct {
	TextOption
	BgHexColor string    `json:"bgHexColor,omitempty"`
	BoxPadding *Padding  `json:"boxPadding,omitempty"`
	BoxSpacing *int      `json:"boxSpacing,omitempty"`
	BoxAlign   box.Align `json:"boxAlign,omitempty"`
	Enabled    *bool     `json:"enabled,omitempty"`
	Limit      int      `json:"limit,omitempty"`
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
