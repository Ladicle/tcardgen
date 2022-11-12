package main

import (
	"log"

	"github.com/spf13/pflag"

	"github.com/hmajid2301/tcardgen/cmd"
)

func init() {
	flags := pflag.NewFlagSet("tcardgen", pflag.ExitOnError)
	pflag.CommandLine = flags
}

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
