package themes

import (
	"cnc/internal"
	"cnc/pkg/logging"
)

var (
	List = map[string]*Theme{
		"default": {
			DisplayName: "Default",
			Description: "The default theme of the command and control.",
			Colors:      []string{"#2a63b8", "#e3e3e3"},
		},
	}
)

type Themes struct{}

func (t *Themes) Serve() {
	err := internal.Config.Unmarshal(&List, "themes")
	if err != nil {
		logging.Global.Fatal().
			Err(err).
			Msg("Failed to unmarshal themes")
	}
}
