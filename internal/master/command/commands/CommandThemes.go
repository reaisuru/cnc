package commands

import (
	"cnc/internal/master/command"
	"cnc/internal/master/models/themes"
	"cnc/internal/master/sessions"
	"cnc/pkg/colorize"
	gradient2 "cnc/pkg/gradient"
	"fmt"
	"github.com/alexeyco/simpletable"
	"strings"
)

func init() {
	var listExecutor = command.ExecutorFunc(func(session *sessions.Session, ctx *command.Context) error {
		var query string
		var err error

		// If the query is not provided, set it to an empty string
		if query, err = ctx.String("query"); err != nil {
			query = ""
		}

		var table = simpletable.New()

		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "#"},
				{Align: simpletable.AlignCenter, Text: "Name"},
				{Align: simpletable.AlignCenter, Text: "Description"},
				{Align: simpletable.AlignCenter, Text: "Colors"},
			},
		}

		table.SetStyle(simpletable.StyleCompact)

		var i int
		for s, theme := range themes.List {
			if !strings.HasPrefix(s, query) {
				continue
			}

			cell := []*simpletable.Cell{
				{Align: simpletable.AlignLeft, Text: fmt.Sprint(i + 1)},
				{Align: simpletable.AlignLeft, Text: theme.DisplayName},
				{Align: simpletable.AlignLeft, Text: theme.Description},
				{Align: simpletable.AlignLeft, Text: generateColorString(theme) + "\x1b[0m"},
			}

			table.Body.Cells = append(table.Body.Cells, cell)
			i++
		}

		return session.Table(table, 1)
	})

	var searchCommand = &command.Command{
		Aliases:     []string{"search", "find"},
		Description: "Searches for a command",
		Arguments: []*command.Argument{
			command.NewArgument("query", nil, command.ArgumentString, true),
		},
		Executor: listExecutor,
	}

	var setCommand = &command.Command{
		Aliases: []string{"set", "use"},
		Arguments: []*command.Argument{
			command.NewArgument("name", nil, command.ArgumentString, true),
		},
		Description: "Sets the theme to use",
		Executor: func(session *sessions.Session, ctx *command.Context) error {
			var name string
			var err error

			// If the query is not provided, set it to an empty string
			if name, err = ctx.String("name"); err != nil {
				return err
			}

			theme, ok := themes.List[strings.ToLower(name)]
			if !ok {
				return session.Println("Theme not found")
			}

			session.Theme = theme
			session.UserProfile.Theme = strings.ToLower(name)

			if err := session.Modify(); err != nil {
				return err
			}

			return session.Printfln("Theme set")
		},
	}

	command.Create(&command.Command{
		Aliases:     []string{"theme", "design", "themes", "designs"},
		Description: "Lists all available themes",
		Roles:       []string{"admin"},
		SubCommands: []*command.Command{
			searchCommand,
			setCommand,
			command.SubCommandListCommand,
		},
		Executor: listExecutor,
	})
}

func generateColorString(theme *themes.Theme) string {
	if theme.IsGradient {
		return gradient2.New(theme.Colors...).Apply(gradient2.Background, strings.Repeat(" ", 14))
	}

	repeatsPerColor := 15 / len(theme.Colors)
	var color string

	for _, hexColor := range theme.Colors {
		color += colorize.Hex(hexColor, true) + strings.Repeat(" ", repeatsPerColor)
	}

	return color + "\x1b[0m"
}
