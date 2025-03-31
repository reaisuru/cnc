package floods

var Presets = map[string]string{
	"valve":   "ffffffff54536f7572636520456e67696e6520517565727900",
	"discord": "1337cafe01000000",
}

func availablePresets() (presets []string) {
	for name, _ := range Presets {
		presets = append(presets, name)
	}

	return presets
}
