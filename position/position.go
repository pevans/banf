package position

import (
	"fmt"
	"strings"
)

func Show(str string, pos int) string {
	// We need to find the beginning and ending of the string we want to
	// highlight. Ideally, this would be the whole line; but we might need to
	// limit that.

	var start, end int

	for start = pos; start > 0 && (pos-30) < start; start-- {
		if str[start] == '\n' {
			break
		}
	}

	for end = pos; end < len(str) && (pos+30) > end; end++ {
		if str[end] == '\n' {
			break
		}
	}

	return fmt.Sprintf(
		"%s\n%s^",
		str[start:end],
		strings.Repeat(" ", (pos-start)),
	)
}
