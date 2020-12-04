package gen

import (
	"bytes"
	"fmt"
	"strings"
)

func (File) LowerName(name string) string {
	return LowerName(name)
}

func (File) UpperName(name string) string {
	return UpperName(name)
}

func (File) JsType(goType string) string {
	return getJsType(goType)
}

func (field MessageField) GoFieldType() string {
	fieldType := getJsType(field.Type)
	if field.Pointer {
		fieldType = "*" + fieldType
	}
	if field.MapKey != "" {
		fieldType = "map[" + field.MapKey + "]" + fieldType
	} else if field.IsArray {
		fieldType = "[]" + fieldType
	}
	return fieldType
	//return fmt.Sprintf("%s %s `json:\"%s\"`", field.Name, fieldType, field.JsName)
}

func (field MessageField) JSFieldType() string {
	var fieldType = field.Type //g.FormatType(field.Type)
	if field.MapKey != "" {
		fieldType = fmt.Sprintf("{ [key:%s]: %s }", field.MapKey, fieldType)
	} else if field.IsArray {
		fieldType = fmt.Sprintf("%s[]", fieldType)
	}
	return fieldType
}

func (c Comment) PrintLeading(prefix string) string {
	var buf bytes.Buffer
	if len(c.LeadingDetached) > 0 {
		buf.WriteString(formatComment(prefix, c.LeadingDetached))
		buf.WriteByte('\n')
		buf.WriteByte('\n')
	}
	if c.Leading != "" {
		buf.WriteString(formatComment(prefix, []string{c.Leading}))
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (c Comment) PrintTrailing(prefix string) string {
	return formatComment(prefix, []string{c.Trailing})
}

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

func getJsType(goType string) string {
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
