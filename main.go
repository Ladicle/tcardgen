package main

import (
	"log"

	"github.com/spf13/pflag"

	"github.com/Ladicle/tcardgen/cmd"
)

func init() {
	flags := pflag.NewFlagSet("kubectl-rolesum", pflag.ExitOnError)
	pflag.CommandLine = flags
}

func main() {
	command := cmd.NewRootCmd()
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}
