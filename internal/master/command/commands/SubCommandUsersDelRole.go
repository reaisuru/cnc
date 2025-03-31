package commands

import (
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
)

var commandUsersDelRole = &command.Command{
	Aliases:     []string{"delrole", "delgroup"},
	Roles:       []string{},
	Description: "Remove a group from a user.",
	Arguments: []*command.Argument{
		command.NewArgument("user", nil, command.ArgumentUser, true),
		command.NewArgument("role", nil, command.ArgumentString, true),
	},
	Executor: func(session *sessions.Session, ctx *command.Context) error {
		user, err := ctx.User("user")
		if err != nil {
			return err
		}

		role, err := ctx.String("role")
		if err != nil {
			return err
		}

		// This is improvable. LOL
		var newRoles = make([]string, 0)
		for _, r := range user.Roles {
			if r != role {
				newRoles = append(newRoles, r)
			}
		}

		user.Roles = newRoles
		if err := user.Modify(); err != nil {
			return err
		}

		return session.Notification("%s is no longer in the %s group.", user.Name, role)
	},
}
