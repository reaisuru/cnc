package views

import (
	"cnc/internal/database"
	"cnc/internal/master/sessions"
	"errors"
	"github.com/dustin/go-humanize/english"
)

func ForcePwChange(session *sessions.Session) error {
	var tries = 0

	if session.Password != database.Hash([]byte("changeme")) {
		return nil
	}

	_ = session.Clear()

	_ = session.Printfln("Changing your password.\r\n")
	_ = session.Printfln("For your security and to protect your account, you need to change your password.")
	_ = session.Printfln("Using the same password as our default can make your account vulnerable.")
	_ = session.Printfln("Please choose a unique password to continue.\r\n")

	pass, err := session.Terminal.ReadPassword("New password: ")
	if err != nil {
		return err
	}

	for tries <= 3 {
		repeat, err := session.Terminal.ReadPassword("Retype new Password: ")
		if err != nil {
			return err
		}

		// not correct
		if pass != repeat {
			_ = session.Notification("Sorry, passwords do not match. "+
				"There %s %s left.",
				english.Plural(3-tries, "is", "are"),
				english.Plural(3-tries, "try", "tries"),
			)
			tries++
			continue
		}

		var hash = database.Hash([]byte(pass))
		if hash == session.Password {
			_ = session.Notification("The password has not been changed. "+
				"There %s %s left.",
				english.Plural(3-tries, "is", "are"),
				english.Plural(3-tries, "try", "tries"),
			)
			tries++
		}

		// modify password in database
		session.Password = hash
		if err := session.Modify(); err != nil {
			return err
		}

		return session.Notification("\x1b[92mpassword updated successfully")
	}

	_ = session.Notification("\x1b[91mFailed preliminary check by password service")
	_, _ = session.Channel.Read(make([]byte, 1))
	return errors.New("failed")
}
