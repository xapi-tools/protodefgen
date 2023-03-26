package protodefgen

import "strings"

type ProtoDefWriter interface {
	ToStringBuilder() (*strings.Builder, error)
	ToFile(string) error
}
