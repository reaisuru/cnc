package colorize

import (
	"fmt"
	"strconv"
	"strings"
)

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func Hex2RGB(hex string) (RGB, error) {
	if strings.HasPrefix(hex, "#") {
		hex = hex[1:]
	}

	var rgb RGB
	values, err := strconv.ParseUint(hex, 16, 32)

	if err != nil {
		return RGB{}, err
	}

	rgb = RGB{
		Red:   uint8(values >> 16),
		Green: uint8((values >> 8) & 0xFF),
		Blue:  uint8(values & 0xFF),
	}

	return rgb, nil
}

func Hex(hex string, background bool) string {
	rgb, err := Hex2RGB(hex)
	if err != nil {
		return ""
	}

	if background {
		return fmt.Sprintf("\u001B[48;2;%d;%d;%dm", rgb.Red, rgb.Green, rgb.Blue)
	}

	return fmt.Sprintf("\u001B[38;2;%d;%d;%dm", rgb.Red, rgb.Green, rgb.Blue)
}
