package commands

import (
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"cnc/pkg/refhelper"
	"strconv"
)

var commandUsersEdit = &command.Command{
	Aliases:     []string{"edit", "modify"},
	Description: "Modifies a user.",
	Arguments: []*command.Argument{
		command.NewArgument("user", nil, command.ArgumentUser, true),
		command.NewArgument("variable", nil, command.ArgumentString, true),
		command.NewArgument("value", nil, command.ArgumentString, true),
	},
	Executor: func(session *sessions.Session, ctx *command.Context) error {
		// the fuckin user
		user, err := ctx.User("user")
		if err != nil {
			return err
		}

		// the variable name which is...obvious I think?
		variable, err := ctx.String("variable")
		if err != nil {
			return err
		}

		// value of the variable ig
		value, err := ctx.String("value")
		if err != nil {
			return err
		}

		// set the variable
		err = refhelper.Set(user, variable, value)
		if err != nil {
			return session.Notification("An error occurred while trying to set the variable: %s", err.Error())
		}

		// modify user in database
		err = user.Modify()
		if err != nil {
			return session.Notification("An error occurred while trying to modify the account: %s", err.Error())
		}

		return session.Notification("Successfully changed %s for %s.", variable, strconv.Quote(user.Name))
	},
}
