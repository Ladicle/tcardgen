package cmd

import (
	"errors"
	"io"
	"os"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/spf13/cobra"
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

	cfm pageparser.ContentFrontMatter
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
	o.cfm, err = pageparser.ParseFrontMatterAndContent(file)
	return err
}

func (o *RootCommandOption) Run(streams IOStreams) error {
	return nil
}
