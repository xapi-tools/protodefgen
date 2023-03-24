package protofilegen

import (
	"fmt"
	"os"
	"strings"
)

type ProtoWriterOpts struct {
	IndentWidth uint
}

type protoFileWriter struct {
	sb    strings.Builder
	Proto *Proto
	Opts  *ProtoWriterOpts
}

func NewProtoFileWriter(p *Proto, o *ProtoWriterOpts) ProtoFileWriter {
	return &protoFileWriter{
		Proto: p,
		Opts:  o,
	}
}

func (pw *protoFileWriter) ToStringBuilder() (*strings.Builder, error) {
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

	if err := pw.addMessages(); err != nil {
		return nil, fmt.Errorf("could not add messages: %v", err)
	}

	if err := pw.writeLine("", 0); err != nil {
		return nil, fmt.Errorf("could not write empty line: %v", err)
	}
	return &pw.sb, nil
}

func (pw *protoFileWriter) ToFile(path string) error {
	sb, err := pw.ToStringBuilder()
	if err != nil {
		return fmt.Errorf("could not get string builder from proto: %v", err)
	}

	if err := os.WriteFile(path, []byte(sb.String()), os.ModePerm); err != nil {
		return fmt.Errorf("could not write proto to file %s: %v", path, err)
	}
	return nil
}

func (pw *protoFileWriter) writeLine(input string, indents uint) error {
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

func (pw *protoFileWriter) addDescription(description string, indents uint) error {
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

func (pw *protoFileWriter) addSyntax() error {
	if err := pw.writeLine("syntax = \"proto3\";", 0); err != nil {
		return fmt.Errorf("could not write syntax: %v", err)
	}

	if err := pw.writeLine("", 0); err != nil {
		return fmt.Errorf("could not write empty line: %v", err)
	}

	return nil
}

func (pw *protoFileWriter) addPackage() error {
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

func (pw *protoFileWriter) addImports() error {
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

func (pw *protoFileWriter) addImport(imp string) error {
	if isEmptyStr(imp) {
		return fmt.Errorf("import cannot be empty")
	}
	if err := pw.writeLine(fmt.Sprintf("import \"%s\";", imp), 0); err != nil {
		return fmt.Errorf("could not add import: %v", err)
	}
	return nil
}

func (pw *protoFileWriter) addMessages() error {
	for i := range pw.Proto.Messages {
		if err := pw.addMessage(&pw.Proto.Messages[i], 0); err != nil {
			return fmt.Errorf("could not add message %s at index %d: %v", pw.Proto.Messages[i].Name, i, err)
		}
	}

	return nil
}

func (pw *protoFileWriter) addMessage(m *Message, indents uint) error {
	if err := pw.addDescription(m.Description, indents); err != nil {
		return fmt.Errorf("could not add description for message: %v", err)
	}

	if isEmptyStr(m.Name) {
		return fmt.Errorf("message name cannot be empty")
	}

	if err := pw.writeLine(fmt.Sprintf("message %s {", m.Name), indents); err != nil {
		return fmt.Errorf("could not add message start: %v", err)
	}

	for i := range m.Fields {
		if err := pw.addField(&m.Fields[i], indents+1); err != nil {
			return fmt.Errorf("could not add field %s at index %d: %v", m.Fields[i].Name, i, err)
		}
	}

	if err := pw.writeLine("}\n", indents); err != nil {
		return fmt.Errorf("could not add message end: %v", err)
	}
	return nil
}

func (pw *protoFileWriter) addField(f *MessageField, indents uint) error {
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