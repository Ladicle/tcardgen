package cmd

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/golang/freetype/truetype"
	"github.com/spf13/cobra"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	tgimg "github.com/Ladicle/tcardgen/pkg/image"
)

var (
	// set values via build flags
	command string
	version string
	commit  string
)

type IOStreams struct {
	Out    io.Writer
	ErrOut io.Writer
}

type RootCommandOption struct {
	Filename   string
	ConfigPath string

	title     string
	author    string
	category  string
	tags      []string
	updatedAt time.Time
}

func NewRootCmd() *cobra.Command {
	opt := RootCommandOption{}
	cmd := &cobra.Command{
		Use:                   "tcardgen [-c <CONFIGURATION>] <FILENAME>",
		Version:               version,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Short:                 "Generate twitter card image from the Hugo post.",
		RunE: func(cmd *cobra.Command, args []string) error {
			streams := IOStreams{
				Out:    os.Stdout,
				ErrOut: os.Stderr,
			}
			if err := opt.Validate(cmd, args); err != nil {
				return err
			}
			if err := opt.Complete(); err != nil {
				return err
			}
			return opt.Run(streams)
		},
	}
	cmd.Flags().StringVarP(&opt.ConfigPath, "config", "c", "config.tcard.yaml", "Set tcardgen configuration path.")
	return cmd
}

func (o *RootCommandOption) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("required argument <FILENAME> is not set")
	}
	o.Filename = args[0]
	return nil
}

func (o *RootCommandOption) Complete() error {
	// TODO: Load configuration

	file, err := os.Open(o.Filename)
	if err != nil {
		return err
	}
	cfm, err := pageparser.ParseFrontMatterAndContent(file)

	title := cfm.FrontMatter["title"].(string)
	if title == "" {
		return fmt.Errorf("can not get title from front matter: %+v", cfm.FrontMatter)
	}
	o.title = title

	if o.author, err = getFirstFMItem(cfm, "author"); err != nil {
		return err
	}

	if o.category, err = getFirstFMItem(cfm, "categories"); err != nil {
		return err
	}

	tags := cfm.FrontMatter["tags"].([]interface{})
	if len(tags) < 1 {
		return fmt.Errorf("can not get tags from front matter: %+v", cfm.FrontMatter)
	}
	for _, t := range tags {
		o.tags = append(o.tags, t.(string))
	}

	o.updatedAt, err = time.Parse("2006-01-02T15:04:05-07:00", cfm.FrontMatter["lastmod"].(string))
	return err
}

func getFirstFMItem(cfm pageparser.ContentFrontMatter, key string) (string, error) {
	categoriesitems := cfm.FrontMatter[key].([]interface{})
	if len(categoriesitems) < 1 {
		return "", fmt.Errorf("can not get %s from front matter: %+v", key, cfm.FrontMatter)
	}
	return categoriesitems[0].(string), nil
}

const (
	fontDir      = "font"
	templateFile = "template.png"
	outputPath   = "thumbnail.png"
)

func (o *RootCommandOption) Run(streams IOStreams) error {
	// load fonts
	ff, err := tgimg.LoadFontFamilyFromDir(fontDir)
	if err != nil {
		return err
	}

	// load template
	tpl, err := tgimg.LoadFromFile(templateFile)
	if err != nil {
		return err
	}

	// write template
	dst := image.NewRGBA(tpl.Bounds())
	draw.Draw(dst, dst.Bounds(), tpl, image.Point{}, draw.Over)

	// write texts
	dr := &font.Drawer{Dst: dst, Dot: fixed.Point26_6{}}

	dr.Face = truetype.NewFace(ff.GetFont(tgimg.FontStyleBold), &truetype.Options{Size: 72})
	dr.Src = image.Black
	dr.Dot.X = fixed.I(127)
	dr.Dot.Y = fixed.I(173 + 72)
	if err := tgimg.DrawText(dr, o.title, 946); err != nil {
		return err
	}

	gray, err := tgimg.Hex("#8D8D8D")
	if err != nil {
		return err
	}
	dr.Face = truetype.NewFace(ff.GetFont(tgimg.FontStyleRegular), &truetype.Options{Size: 42})
	dr.Src = image.NewUniform(gray)
	dr.Dot.X = fixed.I(130)
	dr.Dot.Y = fixed.I(124 + 42)
	if err := tgimg.DrawText(dr, strings.ToUpper(o.category), 946); err != nil {
		return err
	}

	info := fmt.Sprintf("%s・%s", o.author, o.updatedAt.Format("Jan 2"))
	dr.Face = truetype.NewFace(ff.GetFont(tgimg.FontStyleRegular), &truetype.Options{Size: 38})
	dr.Dot.X = fixed.I(231)
	dr.Dot.Y = fixed.I(449 + 38)
	if err := tgimg.DrawText(dr, info, 946); err != nil {
		return err
	}

	tgimg.SaveAsPNG(outputPath, dst)
	return nil
}
