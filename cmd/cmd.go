package cmd

import (
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/spf13/cobra"

	"github.com/Ladicle/tcardgen/pkg/canvas"
	"github.com/Ladicle/tcardgen/pkg/canvas/fontfamily"
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
		o.tags = append(o.tags, strings.Title(t.(string)))
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
	ffa, err := fontfamily.LoadFromDir(fontDir)
	if err != nil {
		return err
	}
	fmt.Fprintf(streams.Out, "Load fonts from %v\n", fontDir)

	c, err := canvas.CreateCanvasFromImage(templateFile)
	if err != nil {
		return err
	}
	fmt.Fprintf(streams.Out, "Load %v template\n", templateFile)

	if err := c.DrawTextAtPoint(
		o.title,
		123, 165,
		canvas.MaxWidth(946),
		canvas.LineSpace(10),
		canvas.FgColor(image.Black),
		canvas.FontFaceFromFFA(ffa, fontfamily.Bold, 72)); err != nil {
		return err
	}
	if err := c.DrawTextAtPoint(
		strings.ToUpper(o.category),
		126, 119,
		canvas.FgHexColor("#8D8D8D"),
		canvas.FontFaceFromFFA(ffa, fontfamily.Regular, 42)); err != nil {
		return err
	}
	if err := c.DrawTextAtPoint(
		fmt.Sprintf("%s・%s", o.author, o.updatedAt.Format("Jan 2")),
		227, 441,
		canvas.FontFaceFromFFA(ffa, fontfamily.Regular, 38)); err != nil {
		return err
	}
	if err := c.DrawBoxTexts(
		o.tags,
		814, 451,
		canvas.FgColor(image.White),
		canvas.BgHexColor("#60BCE0"),
		canvas.BoxPadding(6, 10, 6, 10),
		canvas.BoxSpacing(6),
		canvas.FontFaceFromFFA(ffa, fontfamily.Medium, 22)); err != nil {
		return err
	}

	if err := c.SaveAsPNG(outputPath); err != nil {
		return err
	}
	fmt.Fprintf(streams.Out, "Save image to %v\n", outputPath)
	return nil
}
