package themes

type Theme struct {
	DisplayName string   `swash:"name" toml:"display-name"`
	Description string   `swash:"description" toml:"description"`
	Colors      []string `swash:"colors" toml:"colors"`
	IsGradient  bool     `swash:"gradient" toml:"is-gradient"`
}
