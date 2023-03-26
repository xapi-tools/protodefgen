package protodefgen

import "strings"

func isEmptyStr(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
