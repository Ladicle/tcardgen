package main

import (
	"log"

	"github.com/Ladicle/tcardgen/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
