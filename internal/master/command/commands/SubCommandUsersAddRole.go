package commands

import (
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
)

var commandUsersAddRole = &command.Command{
	Aliases:     []string{"addrole", "addgroup"},
	Roles:       []string{},
	Description: "Add a group to a user.",
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

		// check if user already has the role
		if user.HasRole(role) {
			return session.Notification("%s already has the role.", user.Name)
		}

		// append role to user and
		user.Roles = append(user.Roles, role)
		if err := user.Modify(); err != nil {
			return err
		}

		return session.Notification("%s is now in the %s group.", user.Name, role)
	},
}
