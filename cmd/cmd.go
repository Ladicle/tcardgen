package cmd

import (
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Ladicle/tcardgen/pkg/canvas"
	"github.com/Ladicle/tcardgen/pkg/canvas/fontfamily"
	"github.com/Ladicle/tcardgen/pkg/hugo"
)

const (
	defaultTplImg  = "template.png"
	defaultFontDir = "font"
	defaultOutDir  = "out"

	longDesc = `Generate TwitterCard(OGP) images for your Hugo posts.
Supported front-matters are title, author, categories, tags, and date.`
	example = `# Generate a image and output to the example directory.
tcardgen --fontDir=font --outDir=example --template=example/template.png example/blog-post.md

# Generate multiple images.
tcardgen --template=example/template.png example/*.md`
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
	files   []string
	fontDir string
	outDir  string
	tplImg  string
}

func NewRootCmd() *cobra.Command {
	opt := RootCommandOption{}
	cmd := &cobra.Command{
		Use:                   "tcardgen [-f <FONTDIR>] [-o <OUTDIR>] [-t <TEMPLATE>] <FILE>...",
		Version:               version,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Short:                 "Generate TwitterCard(OGP) image for your Hugo posts.",
		Long:                  longDesc,
		Example:               example,
		RunE: func(cmd *cobra.Command, args []string) error {
			streams := IOStreams{
				Out:    os.Stdout,
				ErrOut: os.Stderr,
			}
			if err := opt.Validate(cmd, args); err != nil {
				return err
			}
			return opt.Run(streams)
		},
	}
	cmd.Flags().StringVarP(&opt.fontDir, "fontDir", "f", defaultFontDir, "Set a font directory.")
	cmd.Flags().StringVarP(&opt.outDir, "outDir", "o", defaultOutDir, "Set an output directory.")
	cmd.Flags().StringVarP(&opt.tplImg, "template", "t", defaultTplImg, "Set a template image file.")
	return cmd
}

func (o *RootCommandOption) Validate(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("required argument <FILE> is not set")
	}
	o.files = args
	return nil
}

func (o *RootCommandOption) Run(streams IOStreams) error {
	ffa, err := fontfamily.LoadFromDir(o.fontDir)
	if err != nil {
		return err
	}
	fmt.Fprintf(streams.Out, "Load fonts from %v\n", o.fontDir)

	tpl, err := canvas.LoadFromFile(o.tplImg)
	if err != nil {
		return err
	}
	fmt.Fprintf(streams.Out, "Load %v template\n", o.tplImg)

	if _, err := os.Stat(o.outDir); os.IsNotExist(err) {
		err := os.Mkdir(o.outDir, 0755)
		if err != nil {
			return err
		}
	}

	var errCnt int
	for _, f := range o.files {
		base := filepath.Base(f)
		name := base[:len(base)-len(filepath.Ext(base))]
		out := filepath.Join(o.outDir, fmt.Sprintf("%s.png", name))

		if err := generateTCard(f, out, tpl, ffa); err != nil {
			fmt.Fprintf(streams.ErrOut, "Failed to generate twitter card for %v: %v\n", out, err)
			errCnt++
			continue
		}
		fmt.Fprintf(streams.Out, "Success to generate twitter card into %v\n", out)
	}

	if errCnt != 0 {
		return fmt.Errorf("failed to generate %d twitter cards", errCnt)
	}
	return nil
}

func generateTCard(contentPath, outPath string, tpl image.Image, ffa *fontfamily.FontFamily) error {
	fm, err := hugo.ParseFrontMatter(contentPath)
	if err != nil {
		return err
	}

	c, err := canvas.CreateCanvasFromImage(tpl)
	if err != nil {
		return err
	}

	var tags []string
	for _, t := range fm.Tags {
		tags = append(tags, strings.Title(t))
	}

	if err := c.DrawTextAtPoint(
		fm.Title,
		123, 165,
		canvas.MaxWidth(946),
		canvas.LineSpace(10),
		canvas.FgColor(image.Black),
		canvas.FontFaceFromFFA(ffa, fontfamily.Bold, 72)); err != nil {
		return err
	}
	if err := c.DrawTextAtPoint(
		strings.ToUpper(fm.Category),
		126, 119,
		canvas.FgHexColor("#8D8D8D"),
		canvas.FontFaceFromFFA(ffa, fontfamily.Regular, 42)); err != nil {
		return err
	}
	if err := c.DrawTextAtPoint(
		fmt.Sprintf("%s・%s", fm.Author, fm.Date.Format("Jan 2")),
		227, 441,
		canvas.FontFaceFromFFA(ffa, fontfamily.Regular, 38)); err != nil {
		return err
	}
	if err := c.DrawBoxTexts(
		tags,
		1025, 451,
		canvas.FgColor(image.White),
		canvas.BgHexColor("#60BCE0"),
		canvas.BoxPadding(6, 10, 6, 10),
		canvas.BoxSpacing(6),
		canvas.BoxAlign(canvas.AlineRight),
		canvas.FontFaceFromFFA(ffa, fontfamily.Medium, 22)); err != nil {
		return err
	}

	return c.SaveAsPNG(outPath)
}
