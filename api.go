package protofilegen

import "strings"

type ProtoFileWriter interface {
	ToStringBuilder() (*strings.Builder, error)
	ToFile(string) error
}
