package internal

import (
	"cnc/pkg/simpleconfig"
	"github.com/alexeyco/simpletable"
	"time"
)

var (
	Started = time.Now()
	Config  = simpleconfig.New(simpleconfig.Toml, "resources/config")

	MainColor = "\x1b[38;2;47;67;168m"

	RawAttacksEnabled = true
	ApiAttacksEnabled = true

	GlobalSlots    = 1
	ApiSlots       = 2
	GlobalCooldown = 15 * time.Second

	StyleSimple = &simpletable.Style{
		Border: &simpletable.BorderStyle{
			TopLeft:            "",
			Top:                "",
			TopRight:           "",
			Right:              "",
			BottomRight:        "",
			Bottom:             "",
			BottomLeft:         "",
			Left:               "",
			TopIntersection:    "",
			BottomIntersection: "",
		},
		Divider: &simpletable.DividerStyle{
			Left:         "",
			Center:       "-",
			Right:        "",
			Intersection: "+",
		},
		Cell: "|",
	}
)
