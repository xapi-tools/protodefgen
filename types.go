package protofilegen

type Proto struct {
	Description string
	Package     string
	Options     []Option
	Imports     []string
	Enums       []Enum
	Messages    []Message
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
