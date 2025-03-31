package commands

import (
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
)

func init() {
	command.Create(&command.Command{
		Aliases:     []string{"exit", "quit"},
		Description: "Disconnects you from.",
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			session.Close()
			return nil
		},
	})
}
