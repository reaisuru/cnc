package commands

import (
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"encoding/hex"
)

func init() {
	command.Create(&command.Command{
		Aliases:     []string{"publickey"},
		Description: "Sets a public key for the user.",
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			pass, err := session.Terminal.ReadPassword("Paste public key: ")
			if err != nil {
				return err
			}

			session.PublicKey = hex.EncodeToString([]byte(pass))
			if err := session.Modify(); err != nil {
				return err
			}

			return session.Notification("Public key has been set.")
		},
	})
}
