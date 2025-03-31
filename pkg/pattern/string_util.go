package pattern

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatCount formats an integer if its -1
func FormatCount(c int, v string) string {
	if c == -1 {
		return v
	}

	return strconv.Itoa(c)
}

// FormatArray formats an array into a PostgreSQL stylish array string.
func FormatArray[T any](v []T) string {
	var str strings.Builder
	str.WriteString("{")

	for i, val := range v {
		if i > 0 {
			str.WriteString(", ")
		}

		str.WriteString(fmt.Sprintf("%v", val))
	}

	str.WriteString("}")
	return str.String()
}
