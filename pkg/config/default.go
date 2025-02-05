package config

import (
	"github.com/Ladicle/tcardgen/pkg/canvas/box"
	"github.com/Ladicle/tcardgen/pkg/canvas/fontfamily"
)

const DefaultTemplate = "example/template.png"

var defaultCnf = DrawingConfig{
	Title: &MultiLineTextOption{
		TextOption: TextOption{
			Start:      &Point{X: 123, Y: 165},
			FgHexColor: "#000000",
			FontSize:   72,
			FontStyle:  fontfamily.Bold,
		},
		MaxWidth:    946,
		LineSpacing: ptrInt(10),
	},
	Category: &TextOption{
		Enabled:    ptrBool(true),
		Start:      &Point{X: 126, Y: 119},
		FgHexColor: "#8D8D8D",
		FontSize:   42,
		FontStyle:  fontfamily.Regular,
	},
	Info: &TextOption{
		Enabled:    ptrBool(true),
		Start:      &Point{X: 227, Y: 441},
		FgHexColor: "#8D8D8D",
		FontSize:   38,
		FontStyle:  fontfamily.Regular,
		Separator:  "・",
		TimeFormat: "Jan 2",
	},
	Tags: &BoxTextsOption{
		Enabled:    ptrBool(true),
		Limit:      0,
		TextOption: TextOption{
			Start:      &Point{X: 1025, Y: 451},
			FgHexColor: "#FFFFFF",
			FontSize:   22,
			FontStyle:  fontfamily.Medium,
		},
		BgHexColor: "#60BCE0",
		BoxPadding: &Padding{Top: 6, Right: 10, Bottom: 6, Left: 10},
		BoxSpacing: ptrInt(6),
		BoxAlign:   box.AlignRight,
	},
}

func Defaulting(cnf *DrawingConfig, tplImg string) {
	if tplImg != "" {
		cnf.Template = tplImg
	} else if cnf.Template == "" {
		cnf.Template = DefaultTemplate
	}

	if cnf.Title == nil {
		cnf.Title = &MultiLineTextOption{}
	}
	defaultingTitle(cnf.Title)

	if cnf.Category == nil {
		cnf.Category = &TextOption{}
	}
	defaultingCategory(cnf.Category)

	if cnf.Info == nil {
		cnf.Info = &TextOption{}
	}
	defaultingInfo(cnf.Info)

	if cnf.Tags == nil {
		cnf.Tags = &BoxTextsOption{}
	}
	defaultTags(cnf.Tags)
}

func defaultingTitle(mto *MultiLineTextOption) {
	setArgsAsDefaultTextOption(&mto.TextOption, &defaultCnf.Title.TextOption)
	if mto.MaxWidth == 0 {
		mto.MaxWidth = defaultCnf.Title.MaxWidth
	}
	if mto.LineSpacing == nil {
		mto.LineSpacing = defaultCnf.Title.LineSpacing
	}
}

func defaultingCategory(to *TextOption) {
	setArgsAsDefaultTextOption(to, defaultCnf.Category)
}

func defaultingInfo(to *TextOption) {
	setArgsAsDefaultTextOption(to, defaultCnf.Info)
}

func defaultTags(bto *BoxTextsOption) {
	if bto == nil {
		bto = &BoxTextsOption{}
	}
	if bto.Enabled == nil {
		bto.Enabled = defaultCnf.Tags.Enabled
	}
	if bto.Limit < 0 {
		bto.Limit = defaultCnf.Tags.Limit
	}

	setArgsAsDefaultTextOption(&bto.TextOption, &defaultCnf.Tags.TextOption)

	if bto.BgHexColor == "" {
		bto.BgHexColor = defaultCnf.Tags.BgHexColor
	}
	if bto.BoxPadding == nil {
		bto.BoxPadding = defaultCnf.Tags.BoxPadding
	}
	if bto.BoxSpacing == nil {
		bto.BoxSpacing = defaultCnf.Tags.BoxSpacing
	}
	if bto.BoxAlign == "" {
		bto.BoxAlign = defaultCnf.Tags.BoxAlign
	}
}

func setArgsAsDefaultTextOption(to *TextOption, dto *TextOption) {
	if to.Enabled == nil {
		to.Enabled = dto.Enabled
	}
	if to.Start == nil {
		to.Start = &Point{X: dto.Start.X, Y: dto.Start.Y}
	}
	if to.FgHexColor == "" {
		to.FgHexColor = dto.FgHexColor
	}
	if to.FontSize == 0 {
		to.FontSize = dto.FontSize
	}
	if to.FontStyle == "" {
		to.FontStyle = dto.FontStyle
	}
	if to.Separator == "" {
		to.Separator = dto.Separator
	}
	if to.TimeFormat == "" {
		to.TimeFormat = dto.TimeFormat
	}
}

func ptrInt(x int) *int {
	return &x
}

func ptrBool(b bool) *bool {
    return &b
}
