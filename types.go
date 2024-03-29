package protodefgen

type Proto struct {
	Description string
	Package     string
	Options     []Option
	Imports     []string
	Enums       []Enum
	Messages    []Message
	Services    []Service
}

type Option struct {
	GoPackage string
}

type Message struct {
	Description string
	Name        string
	Fields      []MessageField
	Enums       []Enum
	Messages    []Message
}

type MessageField struct {
	Description string
	Id          uint
	Name        string
	Type        string
	Optional    bool
	Repeated    bool
}

type Enum struct {
	Description string
	Name        string
	Constants   []EnumConstant
}

type EnumConstant struct {
	Description string
	Name        string
	Value       uint
}

type Service struct {
	Description string
	Name        string
	Methods     []ServiceMethod
}

type ServiceMethod struct {
	Description    string
	Name           string
	Request        string
	StreamRequest  bool
	Response       string
	StreamResponse bool
}
