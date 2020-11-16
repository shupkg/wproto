package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

func New() *Gen {
	return &Gen{
		Out:      ".",
		Root:     "",
		Language: nil,
		Protoc:   "/usr/local/bin/protoc",
	}
}

type Gen struct {
	Out      string   //输出文件夹
	Root     string   //跟路径包名，只有包名以他开头的才会编码
	Language []string //输出语言，目前支持 go 和 ts
	Files    []string //proto 文件
	Protoc   string   //protoc 路径
}

func (g *Gen) Walk() []string {
	var root, _ = filepath.Abs(".")
	var files []string
	for _, f := range g.Files {
		f, _ = filepath.Abs(f)
		_ = filepath.Walk(f, func(truePath string, info os.FileInfo, err error) error {
			if !strings.HasPrefix(truePath, root) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if !info.IsDir() && filepath.Ext(truePath) == ".proto" {
				path, _ := filepath.Rel(root, truePath)
				files = append(files, path)
			}

			return nil
		})
	}
	return files
}

func (g *Gen) genParams() string {
	w := bytes.Buffer{}
	for _, lang := range g.Language {
		lang = strings.TrimFunc(lang, func(r rune) bool { return r == '=' || unicode.IsSpace(r) })
		if !strings.Contains(lang, "=") {
			lang += "=."
		}
		w.WriteString(lang)
		w.WriteByte('&')
	}
	if g.Root != "" {
		w.WriteString("root=")
		w.WriteString(g.Root)
		w.WriteByte('&')
	}

	if w.Len() > 0 {
		w.Truncate(w.Len() - 1)
	}

	w.WriteByte(':')
	if g.Out == "" {
		g.Out = "."
	}
	w.WriteString(g.Out)

	return w.String()
}

func (g *Gen) Exec() error {
	var args []string
	bin, _ := os.Executable()
	args = append(args, "--plugin=protoc-gen-gos="+bin)
	args = append(args, fmt.Sprintf("--gos_out=%s", g.genParams()))
	args = append(args, "--proto_path=.")
	args = append(args, g.Walk()...)
	protoCmd := exec.Command(g.Protoc, args...)
	protoCmd.Stderr = os.Stderr

	log.Println(protoCmd.String())
	log.Println()
	return protoCmd.Run()
}
