package gen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/shupkg/wproto/gen/templates"
)

type TsFormatter struct {
}

func (g *TsFormatter) Template() string {
	return string(templates.TsGots.Bytes())
}

func (g *TsFormatter) FormatImport(i Import, f File) string {
	if i.Path == "google/protobuf/empty.proto" {
		return ""
	}
	path, _ := filepath.Rel(filepath.Dir(f.Path), i.Path)
	path = strings.TrimSuffix(path, ".proto") //+"/"+filepath.Base(i.Path)
	if path[0] != '.' {
		path = "./" + path
	}
	return fmt.Sprintf(`import * as %s from "%s"`, i.Object, path)
}

func (g *TsFormatter) FormatType(goType string) string {
	switch goType {
	case "float32", "float64", "int", "int32", "int64", "uint", "uint32", "uint64":
		return "number"
	case "bool":
		return "boolean"
	case "string":
		return "string"
	case "bytes":
		return "number"
	case "empty.Empty":
		//return "never"
		return "{}"
	default:
		return goType
	}
}

func (g *TsFormatter) FormatField(field MessageField) string {
	var fieldType = g.FormatType(field.Type)
	if field.MapKey != "" {
		fieldType = fmt.Sprintf("{ [key:%s]: %s }", field.MapKey, fieldType)
	} else if field.IsArray {
		fieldType = fmt.Sprintf("%s[]", fieldType)
	}
	return fmt.Sprintf(`%s: %s;`, field.JsName, fieldType)
}

func (g *TsFormatter) Leading(c Comment) string {
	var buf bytes.Buffer
	if len(c.LeadingDetached) > 0 {
		buf.WriteString(formatComment("//", c.LeadingDetached))
		buf.WriteByte('\n')
		buf.WriteByte('\n')
	}
	if c.Leading != "" {
		buf.WriteString(formatComment("//", []string{c.Leading}))
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (g *TsFormatter) Trailing(c Comment) string {
	return formatComment("//", []string{c.Trailing})
}

var _ Formatter = (*TsFormatter)(nil)
