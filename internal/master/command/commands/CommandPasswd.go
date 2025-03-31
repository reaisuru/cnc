package commands

import (
	"cnc/internal/database"
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"github.com/dustin/go-humanize/english"
)

func init() {
	command.Create(&command.Command{
		Aliases:     []string{"passwd"},
		Description: "Change user password.",
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			var tries = 0

			pass, err := session.Terminal.ReadPassword("Password: ")
			if err != nil {
				return err
			}

			for tries <= 3 {
				repeat, err := session.Terminal.ReadPassword("Repeat Password: ")
				if err != nil {
					return err
				}

				// not correct
				if pass != repeat {
					_ = session.Notification("Repeated password is not the same. There is %s left.", english.Plural(3-tries, "try", "tries"))
					tries++
					continue
				}

				// hash password, set password and modify user profile.
				session.Password = database.Hash([]byte(pass))
				if err := session.Modify(); err != nil {
					return err
				}

				return session.Notification("Password successfully changed.")
			}

			return session.Notification("Password change failed. Try again.")
		},
	})
}
