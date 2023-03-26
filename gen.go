package protodefgen

import (
	"fmt"
	"os"
	"strings"
)

type ProtoWriterOpts struct {
	IndentWidth uint
}

type protoDefWriter struct {
	sb    strings.Builder
	Proto *Proto
	Opts  *ProtoWriterOpts
}

func NewProtoDefWriter(p *Proto, o *ProtoWriterOpts) ProtoDefWriter {
	return &protoDefWriter{
		Proto: p,
		Opts:  o,
	}
}

func (pw *protoDefWriter) ToStringBuilder() (*strings.Builder, error) {
	if err := pw.addDescription(pw.Proto.Description, 0); err != nil {
		return nil, fmt.Errorf("could not add description for proto file: %v", err)
	}

	if err := pw.writeLine("", 0); err != nil {
		return nil, fmt.Errorf("could not write empty line: %v", err)
	}

	if err := pw.addSyntax(); err != nil {
		return nil, fmt.Errorf("could not add syntax: %v", err)
	}

	if err := pw.addPackage(); err != nil {
		return nil, fmt.Errorf("could not add package: %v", err)
	}

	if err := pw.addImports(); err != nil {
		return nil, fmt.Errorf("could not add imports: %v", err)
	}

	if err := pw.addEnums(); err != nil {
		return nil, fmt.Errorf("could not add enums: %v", err)
	}

	if err := pw.addMessages(); err != nil {
		return nil, fmt.Errorf("could not add messages: %v", err)
	}

	if err := pw.addServices(); err != nil {
		return nil, fmt.Errorf("could not add services: %v", err)
	}

	return &pw.sb, nil
}

func (pw *protoDefWriter) ToFile(path string) error {
	sb, err := pw.ToStringBuilder()
	if err != nil {
		return fmt.Errorf("could not get string builder from proto: %v", err)
	}

	if err := os.WriteFile(path, []byte(sb.String()), os.ModePerm); err != nil {
		return fmt.Errorf("could not write proto to file %s: %v", path, err)
	}
	return nil
}

func (pw *protoDefWriter) writeLine(input string, indents uint) error {
	// TODO: optimize this
	for i := 0; i < int(indents)*int(pw.Opts.IndentWidth); i++ {
		if _, err := pw.sb.WriteString(" "); err != nil {
			return fmt.Errorf("could not add indent %d: %v", i, err)
		}
	}
	if _, err := pw.sb.WriteString(input); err != nil {
		return fmt.Errorf("could not add input: %v", err)
	}
	if _, err := pw.sb.WriteString("\n"); err != nil {
		return fmt.Errorf("could not add newline: %v", err)
	}
	return nil
}

func (pw *protoDefWriter) addDescription(description string, indents uint) error {
	if isEmptyStr(description) {
		return nil
	}
	for _, line := range strings.Split(description, "\n") {
		if err := pw.writeLine(fmt.Sprintf("// %s", line), indents); err != nil {
			return err
		}
	}

	return nil
}

func (pw *protoDefWriter) addSyntax() error {
	if err := pw.writeLine("syntax = \"proto3\";", 0); err != nil {
		return fmt.Errorf("could not write syntax: %v", err)
	}

	if err := pw.writeLine("", 0); err != nil {
		return fmt.Errorf("could not write empty line: %v", err)
	}

	return nil
}

func (pw *protoDefWriter) addPackage() error {
	if isEmptyStr(pw.Proto.Package) {
		return fmt.Errorf("package name cannot be empty")
	}

	if err := pw.writeLine(fmt.Sprintf("package %s;", pw.Proto.Package), 0); err != nil {
		return fmt.Errorf("could not add package: %v", err)
	}

	if err := pw.writeLine("", 0); err != nil {
		return fmt.Errorf("could not write empty line: %v", err)
	}

	return nil
}

func (pw *protoDefWriter) addImports() error {
	for _, imp := range pw.Proto.Imports {
		if err := pw.addImport(imp); err != nil {
			return fmt.Errorf("could not add import %s: %v", imp, err)
		}
	}

	if len(pw.Proto.Imports) > 0 {
		if err := pw.writeLine("", 0); err != nil {
			return fmt.Errorf("could not write empty line: %v", err)
		}
	}

	return nil
}

func (pw *protoDefWriter) addImport(imp string) error {
	if isEmptyStr(imp) {
		return fmt.Errorf("import cannot be empty")
	}
	if err := pw.writeLine(fmt.Sprintf("import \"%s\";", imp), 0); err != nil {
		return fmt.Errorf("could not add import: %v", err)
	}
	return nil
}

func (pw *protoDefWriter) addMessages() error {
	for i := range pw.Proto.Messages {
		if err := pw.addMessage(&pw.Proto.Messages[i], 0); err != nil {
			return fmt.Errorf("could not add message %s at index %d: %v", pw.Proto.Messages[i].Name, i, err)
		}
	}

	return nil
}

func (pw *protoDefWriter) addMessage(m *Message, indents uint) error {
	if err := pw.addDescription(m.Description, indents); err != nil {
		return fmt.Errorf("could not add description for message: %v", err)
	}

	if isEmptyStr(m.Name) {
		return fmt.Errorf("message name cannot be empty")
	}

	if err := pw.writeLine(fmt.Sprintf("message %s {", m.Name), indents); err != nil {
		return fmt.Errorf("could not add message start: %v", err)
	}

	for i := range m.Enums {
		if err := pw.addEnum(&m.Enums[i], indents+1); err != nil {
			return fmt.Errorf("could not add enum %s at index %d: %v", m.Enums[i].Name, i, err)
		}
	}

	for i := range m.Messages {
		if err := pw.addMessage(&m.Messages[i], indents+1); err != nil {
			return fmt.Errorf("could not add message %s at index %d: %v", m.Messages[i].Name, i, err)
		}
	}

	for i := range m.Fields {
		if err := pw.addMessageField(&m.Fields[i], indents+1); err != nil {
			return fmt.Errorf("could not add message field %s at index %d: %v", m.Fields[i].Name, i, err)
		}
	}

	if err := pw.writeLine("}\n", indents); err != nil {
		return fmt.Errorf("could not add message end: %v", err)
	}
	return nil
}

func (pw *protoDefWriter) addMessageField(f *MessageField, indents uint) error {
	if err := pw.addDescription(f.Description, indents); err != nil {
		return fmt.Errorf("could not add description for field: %v", err)
	}

	if isEmptyStr(f.Name) {
		return fmt.Errorf("field name cannot be empty")
	}
	// TODO: validate type ?
	if isEmptyStr(f.Type) {
		return fmt.Errorf("field type cannot be empty")
	}

	if f.Id < 1 {
		return fmt.Errorf("field id cannot be less than 1")
	}

	cardinality := ""

	if f.Optional {
		cardinality = "optional "
	} else if f.Repeated {
		cardinality = "repeated "
	}

	if err := pw.writeLine(fmt.Sprintf("%s%s %s = %d;", cardinality, f.Type, f.Name, f.Id), indents); err != nil {
		return fmt.Errorf("could not add field line: %v", err)
	}

	return nil
}

func (pw *protoDefWriter) addEnums() error {
	for i := range pw.Proto.Enums {
		if err := pw.addEnum(&pw.Proto.Enums[i], 0); err != nil {
			return fmt.Errorf("could not add enum %s at index %d: %v", pw.Proto.Enums[i].Name, i, err)
		}
	}

	return nil
}

func (pw *protoDefWriter) addEnum(e *Enum, indents uint) error {
	if err := pw.addDescription(e.Description, indents); err != nil {
		return fmt.Errorf("could not add description for enum: %v", err)
	}

	if isEmptyStr(e.Name) {
		return fmt.Errorf("enum name cannot be empty")
	}

	if err := pw.writeLine(fmt.Sprintf("enum %s {", e.Name), indents); err != nil {
		return fmt.Errorf("could not add enum start: %v", err)
	}

	if len(e.Constants) < 1 {
		return fmt.Errorf("must have at least one enum constant")
	}

	if e.Constants[0].Value != 0 {
		return fmt.Errorf("first enum constant must have a value 0")
	}

	for i := range e.Constants {
		if err := pw.addEnumConstant(&e.Constants[i], indents+1); err != nil {
			return fmt.Errorf("could not add enum constant %s at index %d: %v", e.Constants[i].Name, i, err)
		}
	}

	if err := pw.writeLine("}\n", indents); err != nil {
		return fmt.Errorf("could not add enum end: %v", err)
	}
	return nil
}

func (pw *protoDefWriter) addEnumConstant(c *EnumConstant, indents uint) error {
	if err := pw.addDescription(c.Description, indents); err != nil {
		return fmt.Errorf("could not add description for costant: %v", err)
	}

	if isEmptyStr(c.Name) {
		return fmt.Errorf("constant name cannot be empty")
	}

	if err := pw.writeLine(fmt.Sprintf("%s = %d;", c.Name, c.Value), indents); err != nil {
		return fmt.Errorf("could not add constant line: %v", err)
	}

	return nil
}

func (pw *protoDefWriter) addServices() error {
	for i := range pw.Proto.Services {
		if err := pw.addService(&pw.Proto.Services[i], 0); err != nil {
			return fmt.Errorf("could not add service %s at index %d: %v", pw.Proto.Services[i].Name, i, err)
		}
	}

	return nil
}

func (pw *protoDefWriter) addService(s *Service, indents uint) error {
	if err := pw.addDescription(s.Description, indents); err != nil {
		return fmt.Errorf("could not add description for service: %v", err)
	}

	if isEmptyStr(s.Name) {
		return fmt.Errorf("service name cannot be empty")
	}

	if err := pw.writeLine(fmt.Sprintf("service %s {", s.Name), indents); err != nil {
		return fmt.Errorf("could not add service start: %v", err)
	}

	for i := range s.Methods {
		if err := pw.addServiceMethod(&s.Methods[i], indents+1); err != nil {
			return fmt.Errorf("could not add service method %s at index %d: %v", s.Methods[i].Name, i, err)
		}
	}

	if err := pw.writeLine("}\n", indents); err != nil {
		return fmt.Errorf("could not add service end: %v", err)
	}
	return nil
}

func (pw *protoDefWriter) addServiceMethod(m *ServiceMethod, indents uint) error {
	if err := pw.addDescription(m.Description, indents); err != nil {
		return fmt.Errorf("could not add description for method: %v", err)
	}

	if isEmptyStr(m.Name) {
		return fmt.Errorf("method name cannot be empty")
	}
	if isEmptyStr(m.Request) {
		return fmt.Errorf("method request cannot be empty")
	}
	if isEmptyStr(m.Response) {
		return fmt.Errorf("method response cannot be empty")
	}

	streamReq := ""
	if m.StreamRequest {
		streamReq = "stream "
	}

	streamRes := ""
	if m.StreamResponse {
		streamRes = "stream "
	}

	if err := pw.writeLine(fmt.Sprintf("rpc %s(%s%s) returns (%s%s);", m.Name, streamReq, m.Request, streamRes, m.Response), indents); err != nil {
		return fmt.Errorf("could not add method line: %v", err)
	}

	return nil
}
