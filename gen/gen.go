package gen

import (
	"bytes"
	"log"
	"strings"
	"text/template"

	"github.com/shupkg/wproto/gen/templates"

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

func Run(plugin *protogen.Plugin) error {
	//log.Println("parameter:" + plugin.Request.GetParameter())
	params := ParseParam(plugin.Request.GetParameter())
	mod := params.Get("module")

	for _, f := range plugin.Files {
		if mod != "" {
			if !strings.HasPrefix(string(f.GoImportPath), mod) {
				continue
			}
		}
		model := ParseFile(f)
		model.ApiPrefix = strings.Trim(strings.TrimPrefix(model.ImportPath, mod), "/")

		if params.GetBool("gos_model") {
			if len(model.Messages) > 0 || len(model.Enums) > 0 {
				var buf = &bytes.Buffer{}
				err := template.Must(template.New("").Parse(string(templates.GosModelGogo.Bytes()))).Execute(buf, model)
				if err != nil {
					log.Fatal(err)
				}
				//log.Println(f.GeneratedFilenamePrefix + ".gos.model.go")
				plugin.NewGeneratedFile(f.GeneratedFilenamePrefix+".gos.model.go", EmptyImport).Write(GoFmt(buf.Bytes()))
			}
		}

		if params.GetBool("gos_rpc") {
			if len(model.Services) > 0 {
				var buf = &bytes.Buffer{}
				err := template.Must(template.New("").Parse(string(templates.GosRpcGogo.Bytes()))).Execute(buf, model)
				if err != nil {
					log.Fatal(err)
				}
				//log.Println(f.GeneratedFilenamePrefix + ".gos.rpc.go")
				plugin.NewGeneratedFile(f.GeneratedFilenamePrefix+".gos.rpc.go", EmptyImport).Write(GoFmt(buf.Bytes()))
			}
		}

		if params.GetBool("gos_ts") {
			if len(model.Services) > 0 || len(model.Messages) > 0 || len(model.Enums) > 0 {
				var buf = &bytes.Buffer{}
				err := template.Must(template.New("").Parse(string(templates.GosClientGots.Bytes()))).Execute(buf, model)
				if err != nil {
					log.Fatal(err)
				}
				//log.Println(f.GeneratedFilenamePrefix + ".gos.client.ts")
				plugin.NewGeneratedFile(f.GeneratedFilenamePrefix+".gos.client.ts", EmptyImport).Write(buf.Bytes())
			}
		}
	}
	return nil
}
