package global

import (
	"strings"
)

func MergeKey(elements ...string) string {
	if len(elements) == 0 {
		return ""
	}
	return strings.Join(elements, ":")
}
