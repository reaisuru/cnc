package main

import (
	"cnc/internal/api"
	"cnc/internal/clients"
	"cnc/internal/database"
	"cnc/internal/master"
	"cnc/internal/master/models"
	"cnc/internal/master/models/roles"
	"cnc/internal/master/models/themes"
	"cnc/pkg/location"
	"cnc/pkg/logging"
)

var (
	Models = []models.Model{
		new(roles.Roles),
		new(themes.Themes),
	}
)

func main() {
	for _, model := range Models {
		model.Serve()
	}

	if err := location.Load("resources/ipinfo.csv"); err != nil {
		logging.Global.Warn().Msg("Failed to load ip information!")
	}

	database.Serve()
	go clients.Listen()
	go api.Serve()

	if masterListener, err := master.NewListener(1337, "resources/ssh.ppk"); err == nil {
		masterListener.Listen()
	}
}
