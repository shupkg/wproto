package gen

import (
	"fmt"
	"sort"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func ParseFile(gf *protogen.File) File {
	var f = File{}
	f.Path = gf.Desc.Path()
	f.Package = string(gf.GoPackageName)
	f.ImportPath = string(gf.GoImportPath)
	f.ImportMap = map[string]string{}

	for _, gm := range gf.Messages {
		f.Messages = append(f.Messages, f.parseMessage(gm))
	}

	for _, ge := range gf.Enums {
		f.Enums = append(f.Enums, f.parseEnum(ge))
	}

	for _, gs := range gf.Services {
		f.Services = append(f.Services, f.parseService(gs))
	}

	sort.Sort(sortImport(f.Imports))
	return f
}

func (f *File) parseService(gs *protogen.Service) Service {
	var s = Service{}
	s.Name = gs.GoName
	s.Comment = f.parseComment(gs.Comments)
	for _, method := range gs.Methods {
		s.Methods = append(s.Methods, f.parseServiceMethod(method))
	}
	return s
}

func (f *File) parseServiceMethod(gm *protogen.Method) ServiceMethod {
	var s = ServiceMethod{}
	s.Name = gm.GoName
	s.Comment = f.parseComment(gm.Comments)
	s.Input = f.getMessageType(gm.Input)
	s.Output = f.getMessageType(gm.Output)
	return s
}

func (f *File) parseMessage(gm *protogen.Message) Message {
	var s Message
	s.Name = f.getMessageType(gm)
	s.Comment = f.parseComment(gm.Comments)

	for _, e := range gm.Enums {
		f.Enums = append(f.Enums, f.parseEnum(e))
	}

	for _, m := range gm.Messages {
		if !m.Desc.IsMapEntry() {
			f.Messages = append(f.Messages, f.parseMessage(m))
		}
	}

	for _, gField := range gm.Fields {
		s.Fields = append(s.Fields, f.parseMessageField(gField))
	}

	return s
}

func (f *File) parseMessageField(gField *protogen.Field) MessageField {
	var sField MessageField
	sField.Name = string(gField.Desc.Name())
	sField.IsArray = gField.Desc.IsList()
	sField.Index = gField.Desc.Index()
	sField.Comment = f.parseComment(gField.Comments)

	var tagMap map[string]string
	var ss = strings.SplitN(sField.Comment.Trailing, "//", 2)
	for _, c := range ss {
		c = strings.TrimSpace(c)
		if IsWrap(c) {
			tagMap = ParseTags(c)
		} else {
			sField.Comment.Trailing = c
		}
	}

	sField.JsName = tagMap["name"]
	if sField.JsName == "" {
		sField.JsName = SnakeCase(sField.Name)
	}

	//基础类型
	if t, ok := f.getKindType(gField.Desc.Kind()); ok {
		sField.Type = t
		return sField
	}

	//枚举
	if gField.Enum != nil {
		sField.Type = f.getEnumType(gField.Enum)
		return sField
	}

	if gField.Message != nil {
		if gField.Message.Desc.IsMapEntry() { //map
			for _, field := range gField.Message.Fields {
				item := f.parseMessageField(field)
				switch item.Name {
				case "key":
					sField.MapKey = item.Type
				case "value":
					sField.Type = item.Type
				}
			}
			return sField
		}
		sField.Pointer = true
		sField.Type = f.getMessageType(gField.Message)
		return sField
	}

	//sField.Type = f.getMessageType(gField)
	return sField
}

func (f *File) parseEnum(ge *protogen.Enum) Enum {
	var s = Enum{}
	s.Name = f.getEnumType(ge)
	s.Comment = f.parseComment(ge.Comments)
	for _, gVal := range ge.Values {
		s.Fields = append(s.Fields, EnumField{
			Name:    f.getEnumType(ge) + string(gVal.Desc.Name()),
			Value:   string(gVal.Desc.Name()),
			Index:   gVal.Desc.Index(),
			Comment: f.parseComment(gVal.Comments),
		})
	}
	return s
}

func (f *File) parseComment(gc protogen.CommentSet) Comment {
	var c = Comment{}
	c.Leading = strings.TrimSpace(string(gc.Leading))
	c.Trailing = strings.TrimSpace(string(gc.Trailing))
	for _, detached := range gc.LeadingDetached {
		c.LeadingDetached = append(c.LeadingDetached, strings.TrimSpace(string(detached)))
	}
	return c
}

func (f *File) getKindType(kind protoreflect.Kind) (string, bool) {
	switch kind {
	case protoreflect.DoubleKind:
		return "float64", true
	case protoreflect.FloatKind:
		return "float32", true
	case protoreflect.Int64Kind:
		return "int64", true
	case protoreflect.Uint64Kind:
		return "uint64", true
	case protoreflect.Int32Kind:
		return "int", true
	case protoreflect.Fixed64Kind:
		return "uint64", true
	case protoreflect.Fixed32Kind:
		return "uint32", true
	case protoreflect.BoolKind:
		return "bool", true
	case protoreflect.StringKind:
		return "string", true
	case protoreflect.GroupKind:
		return "group", false
	case protoreflect.MessageKind:
		return "message", false
	case protoreflect.BytesKind:
		return "bytes", false
	case protoreflect.Uint32Kind:
		return "uint", true
	case protoreflect.EnumKind:
		return "enum", false
	case protoreflect.Sfixed32Kind:
		return "uint32", true
	case protoreflect.Sfixed64Kind:
		return "uint64", true
	case protoreflect.Sint32Kind:
		return "int32", true
	case protoreflect.Sint64Kind:
		return "int64", true
	default:
		return "", false
	}
}

func (f *File) getEnumType(enum *protogen.Enum) string {
	return f.parseIdent(enum.GoIdent, enum.Location.SourceFile)
}

func (f *File) getMessageType(message *protogen.Message) string {
	return f.parseIdent(message.GoIdent, message.Location.SourceFile)
}

func (f *File) parseIdent(ident protogen.GoIdent, fullName string) string {
	importPath := string(ident.GoImportPath)
	if importPath == f.ImportPath {
		return ident.GoName
	}

	cname := importPath
	if idx := strings.LastIndex(cname, "/"); idx != -1 {
		cname = cname[idx+1:]
	}

	if _, find := f.ImportMap[importPath]; !find {
		f.ImportMap[importPath] = cname
		f.Imports = append(f.Imports, Import{
			Object:  cname,
			Package: importPath,
			Path:    fullName,
		})
	}

	return fmt.Sprintf("%s.%s", cname, ident.GoName)
}

type sortImport []Import

func (s sortImport) Len() int {
	return len(s)
}

func (s sortImport) Less(i, j int) bool {
	return s[i].Package < s[j].Package
}

func (s sortImport) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var _ sort.Interface = (sortImport)(nil)
