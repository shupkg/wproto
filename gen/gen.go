package gen

import (
	"bytes"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

const EmptyImport = protogen.GoImportPath("")

/*
"float32"
"float64"
"int"
"int32"
"int64"
"uint"
"uint32"
"uint64"
"bool"
"string"
"bytes"
*/
type Formatter interface {
	Template() string
	FormatField(field MessageField) string
	Leading(c Comment) string
	Trailing(c Comment) string
	FormatType(goType string) string
	FormatImport(i Import, f File) string
}

func ToFuncMap(f Formatter) template.FuncMap {
	return template.FuncMap{
		"FormatImport":  f.FormatImport,
		"FormatField":   f.FormatField,
		"Leading":       f.Leading,
		"Trailing":      f.Trailing,
		"FormatType":    f.FormatType,
		"UnderlineCase": SnakeCase,
	}
}

func Run(plugin *protogen.Plugin) error {
	params, _ := url.ParseQuery(plugin.Request.GetParameter())
	root := params.Get("root")
	if root != "" {
		if !strings.HasSuffix(root, "/") {
			root += "/"
		}
	}

	for _, f := range plugin.Files {
		//确定文件名
		filename := filepath.Join(
			string(f.GoImportPath),
			strings.TrimSuffix(filepath.Base(f.Desc.Path()), filepath.Ext(f.Desc.Path())),
		)
		if !strings.HasPrefix(filename, root) {
			continue
		}
		filename = strings.TrimPrefix(filename, root)

		model := ParseFile(f)
		model.ApiPrefix = strings.TrimPrefix(model.ImportPath, root)

		goOut := params.Get("go")
		if goOut != "" {
			var (
				formatter = &GoFormatter{}
				buf       = &bytes.Buffer{}
			)

			err := template.Must(template.New("").Funcs(ToFuncMap(formatter)).Parse(formatter.Template())).Execute(buf, model)
			if err != nil {
				log.Fatal(err)
			}
			plugin.NewGeneratedFile(filepath.Join(goOut, filename+".go"), EmptyImport).Write(GoFmt(buf.Bytes()))
		}

		tsOut := params.Get("ts")
		if tsOut != "" {
			var (
				formatter = &TsFormatter{}
				buf       = &bytes.Buffer{}
			)

			err := template.Must(template.New("").Funcs(ToFuncMap(formatter)).Parse(formatter.Template())).Execute(buf, model)
			if err != nil {
				log.Fatal(err)
			}
			plugin.NewGeneratedFile(filepath.Join(tsOut, filename+".ts"), EmptyImport).Write(buf.Bytes())
		}
	}
	return nil
}
