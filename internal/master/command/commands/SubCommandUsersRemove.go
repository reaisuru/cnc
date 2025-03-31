package commands

import (
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"strconv"
)

var commandUsersRemove = &command.Command{
	Aliases:     []string{"remove", "drop", "delete"},
	Roles:       []string{},
	Description: "Delete a user from the database.",
	Arguments: []*command.Argument{
		command.NewArgument("user", nil, command.ArgumentUser, true),
	},
	Executor: func(session *sessions.Session, ctx *command.Context) error {
		user, err := ctx.User("user")
		if err != nil {
			return err
		}

		if err := user.Drop(); err != nil {
			return err
		}

		return session.Notification("Deleted %s from the database.", strconv.Quote(user.Name))
	},
}
