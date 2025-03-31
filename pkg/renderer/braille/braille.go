package braille

const (
	brailleStart = 0x2800
	size         = 8
)

var brailleOrder = [size]int{
	1, 4,
	2, 5,
	3, 6,
	7, 8,
}

func Get(dots [size]bool) string {
	offset := 0

	for idx, bit := range brailleOrder {
		if dots[idx] == true {
			offset |= 1 << (bit - 1)
		}
	}

	return string(rune(brailleStart + offset))
}
