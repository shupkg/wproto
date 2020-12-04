package main

import (
	"log"
	"os"

	"github.com/shupkg/wproto/cmd"
	"github.com/shupkg/wproto/gen"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)

	if cmd.InProtoc() {
		protogen.Options{}.Run(gen.Run)
		return
	}

	c := cmd.WithFlag(pflag.CommandLine)
	pflag.Parse()
	c.Files = append(c.Files, pflag.Args()...)

	if err := c.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
