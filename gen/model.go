package gen

type File struct {
	Path       string            `json:",omitempty"`
	Package    string            `json:",omitempty"`
	ApiPrefix  string            `json:",omitempty"`
	ImportPath string            `json:",omitempty"`
	ImportMap  map[string]string `json:",omitempty"`
	Imports    []Import          `json:",omitempty"`
	Enums      []Enum            `json:",omitempty"`
	Messages   []Message         `json:",omitempty"`
	Services   []Service         `json:",omitempty"`
}

type Import struct {
	Object  string //引用的对象 model
	Package string //引用包名
	Path    string //引用包名
}

type Message struct {
	Name   string
	Fields []MessageField `json:",omitempty"`
	Comment
}

type MessageField struct {
	Name    string
	JsName  string
	Type    string `json:",omitempty"`
	IsArray bool   `json:",omitempty"`
	MapKey  string `json:",omitempty"`
	Pointer bool   `json:",omitempty"`
	Index   int
	Comment
}

type Enum struct {
	Name   string
	Fields []EnumField `json:",omitempty"`
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
	Methods []ServiceMethod `json:",omitempty"`
	Comment
}

type ServiceMethod struct {
	Name   string
	Input  string `json:",omitempty"`
	Output string `json:",omitempty"`
	Comment
}

//注释
type Comment struct {
	LeadingDetached []string `json:",omitempty"`
	Leading         string   `json:",omitempty"`
	Trailing        string   `json:",omitempty"`
}
