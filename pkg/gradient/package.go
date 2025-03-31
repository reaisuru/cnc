package gradient

type SwashGradient struct {
	gradient *Gradient
}

// New will create the new Gradient object
func New(colours ...string) *SwashGradient {
	glamour := NewDerivative()

	for _, c := range colours {
		rgb := hex2rgb(c)
		glamour.AppendRgbToGradient(rgb.Red, rgb.Green, rgb.Blue)
	}

	return &SwashGradient{glamour}
}

func (g *SwashGradient) Apply(mode int, content string) string {
	return g.gradient.Marshal(mode, content)
}

func Fast(content string, mode int, colours ...string) string {
	return New(colours...).Apply(mode, content)
}
