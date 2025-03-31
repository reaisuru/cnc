package commands

import (
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"github.com/dustin/go-humanize"
)

var commandUsersAddExpiry = &command.Command{
	Aliases:     []string{"addexpiry"},
	Description: "Extends the subscription time.",
	Arguments: []*command.Argument{
		command.NewArgument("user", nil, command.ArgumentUser, true),
		command.NewArgument("time", nil, command.ArgumentTime, true),
	},
	Executor: func(session *sessions.Session, ctx *command.Context) error {
		user, err := ctx.User("user")
		if err != nil {
			return err
		}

		expiryTime, err := ctx.Time("time")
		if err != nil {
			return err
		}

		user.Expiry = user.Expiry.Add(expiryTime)
		if err := user.Modify(); err != nil {
			return err
		}

		return session.Notification("The account will now expire in %s.", humanize.Time(user.Expiry))
	},
}
