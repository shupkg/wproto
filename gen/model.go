package gen

type File struct {
	Path       string
	Package    string
	ApiPrefix  string
	ImportPath string
	ImportMap  map[string]string
	Imports    []Import
	Enums      []Enum
	Messages   []Message
	Services   []Service
}

type Import struct {
	Object  string //引用的对象 model
	Package string //引用包名
	Path    string //引用包名
}

type Message struct {
	Name   string
	Fields []MessageField
	Comment
}

type MessageField struct {
	Name    string
	JsName  string
	Type    string
	IsArray bool
	MapKey  string
	Pointer bool
	Index   int
	Comment
}

type Enum struct {
	Name   string
	Fields []EnumField
	Comment
}

type EnumField struct {
	Name  string
	Value string
	Index int
	Comment
}

type Service struct {
	Name    string
	Methods []ServiceMethod
	Comment
}

type ServiceMethod struct {
	Name   string
	Input  string
	Output string
	Comment
}

//注释
type Comment struct {
	LeadingDetached []string
	Leading         string
	Trailing        string
}
