package gen

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/shupkg/wproto/gen/templates"
)

type GoFormatter struct {
}

func (g *GoFormatter) Template() string {
	return string(templates.GoGogo.Bytes())
}

func (g *GoFormatter) FormatImport(i Import, _ File) string {
	return fmt.Sprintf(`"%s"`, i.Package)
}

func (g *GoFormatter) FormatType(goType string) string {
	return goType
}

func (g *GoFormatter) FormatField(field MessageField) string {
	fieldType := field.Type
	if field.Pointer {
		fieldType = "*" + fieldType
	}
	if field.MapKey != "" {
		fieldType = "map[" + field.MapKey + "]" + fieldType
	} else if field.IsArray {
		fieldType = "[]" + fieldType
	}

	return fmt.Sprintf("%s %s `json:\"%s\"`", field.Name, fieldType, field.JsName)
}

func (g *GoFormatter) Leading(c Comment) string {
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

func (g *GoFormatter) Trailing(c Comment) string {
	return formatComment("//", []string{c.Trailing})
}

var _ Formatter = (*GoFormatter)(nil)

func formatComment(prefix string, comments []string) string {
	w := &bytes.Buffer{}
	for _, comment := range comments {
		cs := strings.Split(comment, "\n")
		for _, c := range cs {
			c := strings.TrimSpace(c)
			if c != "" {
				w.WriteString(prefix)
				w.WriteString(c)
				w.WriteString("\n")
			}
		}
	}
	return strings.TrimSpace(w.String())
}
