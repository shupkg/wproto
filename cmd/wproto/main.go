package main

import (
	"errors"
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

	stat, err := os.Stdin.Stat()
	if len(os.Args) == 1 && err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		protogen.Options{}.Run(gen.Run)
		return
	}

	var g = cmd.New()
	pflag.ErrHelp = errors.New("")
	pflag.CommandLine.SortFlags = false
	pflag.StringVarP(&g.Out, "out", "o", g.Out, "输出目录")
	pflag.StringVarP(&g.Root, "root", "r", g.Root, "根包名")
	pflag.StringVar(&g.Protoc, "protoc", g.Protoc, "protoc 路径")
	pflag.StringSliceVar(&g.Files, "files", g.Files, "proto文件(夹)")
	pflag.StringSliceVarP(&g.Language, "language", "l", g.Language, "生成语言, 格式为语言=输出路径, 如: --language go=.")
	pflag.Parse()
	g.Files = append(g.Files, pflag.Args()...)
	if err := g.Exec(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
