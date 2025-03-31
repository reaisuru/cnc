package gradient

import (
	"strconv"
	"strings"
)

type color struct {
	Red   int
	Green int
	Blue  int
}

func hex2rgb(hex string) color {
	if strings.HasPrefix(hex, "#") {
		hex = hex[1:]
	}

	values, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return color{255.0, 255.0, 255.0}
	}

	return color{
		int(values >> 16),
		int((values >> 8) & 0xFF),
		int(values & 0xFF),
	}
}
