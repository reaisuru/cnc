package command

import (
	"cnc/internal/database"
	"cnc/internal/master/sessions"
	"cnc/pkg/pattern"
	"errors"
	"fmt"
	"github.com/alexeyco/simpletable"
	"golang.org/x/exp/slices"
	"log"
	"strings"
)

var (
	ErrCommandAlreadyRegistered    = errors.New("command already registered")
	ErrCommandNotRegistered        = errors.New("command not registered")
	ErrCommandNotEnoughPermissions = errors.New("not enough permissions to run command")

	ErrArgumentNotRegistered = errors.New("argument not registered")
	ErrArgumentInvalidType   = errors.New("tried to get argument of invalid type")
	ErrArgumentRequired      = errors.New("missing required argument")

	SubCommandListCommand = &Command{
		Aliases:     []string{"help", "?"},
		Description: "Provides a list of all subcommands.",
		Executor: func(session *sessions.Session, ctx *Context) error {
			var table = simpletable.New()

			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "#"},
					{Align: simpletable.AlignCenter, Text: "Aliases"},
					{Align: simpletable.AlignCenter, Text: "Description"},
					{Align: simpletable.AlignCenter, Text: "Roles"},
					{Align: simpletable.AlignCenter, Text: "Syntax"},
				},
			}

			for i, cmd := range ctx.Parent().SubCommands {
				if cmd.Aliases[0] == "help" {
					continue
				}

				cell := []*simpletable.Cell{
					{Align: simpletable.AlignLeft, Text: fmt.Sprint(i)},
					{Align: simpletable.AlignLeft, Text: strings.Join(cmd.Aliases, ", ")},
					{Align: simpletable.AlignLeft, Text: cmd.Description},
					{Align: simpletable.AlignLeft, Text: pattern.FormatArray(cmd.Roles)},
					{Align: simpletable.AlignLeft, Text: cmd.PrettyArguments()},
				}

				table.Body.Cells = append(table.Body.Cells, cell)
			}

			table.SetStyle(simpletable.StyleCompact)
			return session.Table(table, 1)
		},
	}

	Commands = make([]*Command, 0)
)

// ExecutorFunc is the function for command
type ExecutorFunc func(session *sessions.Session, ctx *Context) error

// ErrorFunc is the function for errors to execute branding...etc
type ErrorFunc func(session *sessions.Session) error

// Command is a simple command
type Command struct {
	Aliases     []string
	Roles       []string
	Description string
	Arguments   []*Argument
	SubCommands []*Command

	// Executor is the executor of the command
	Executor ExecutorFunc
}

// Create adds a command to the registry
func Create(command *Command) *Command {
	if ByName(Commands, command.Aliases[0]) != nil {
		log.Println("Failed to register command: ", ErrCommandAlreadyRegistered)
		return nil
	}

	Commands = append(Commands, command)
	return command
}

// ByName gets a command by its name
func ByName(list []*Command, value string) *Command {
	for _, command := range list {
		if slices.Contains(command.Aliases, value) {
			return command
		}
	}

	return nil
}

// Parse will try to get a command within an argument array. (Parses sub-command too)
func Parse(profile *database.UserProfile, args ...string) (parent *Command, command *Command, index int, err error) {
	// find parent command based on the first argument lol
	parent = ByName(Commands, args[0])
	if parent == nil {
		return nil, nil, 0, ErrCommandNotRegistered
	}

	// checks if user has enough permissions for parent command
	if !profile.ContainsRole(parent.Roles) {
		return nil, nil, 0, ErrCommandNotEnoughPermissions
	}

	// check if parent command has subcommands
	if len(parent.SubCommands) == 0 || len(args) == 1 {
		return parent, parent, 1, nil
	}

	var actualParent = parent

	for pos, arg := range args[1:] {
		// find subcommand based on second argument...
		child := ByName(parent.SubCommands, arg)
		if child == nil {
			return parent, parent, pos + 1, nil
		}

		// checks if user has enough permissions for sub command
		if !profile.ContainsRole(child.Roles) {
			return nil, nil, pos, ErrCommandNotEnoughPermissions
		}

		parent = child
	}

	// return child command and index
	return actualParent, parent, 2, nil
}

func (c *Command) PrettyArguments() string {
	var argStr strings.Builder
	for _, a := range c.Arguments {
		if a.Required {
			argStr.WriteString(fmt.Sprintf("<%s> ", a.Name))
			continue
		}

		argStr.WriteString(fmt.Sprintf("[%s] ", a.Name))
	}

	return argStr.String()
}

func Names(commands []*Command, p *database.UserProfile) []string {
	var names = make([]string, 0)

	for _, c := range commands {
		if c.Aliases[0] == "help" || !p.ContainsRole(c.Roles) {
			continue
		}

		names = append(names, c.Aliases[0])
	}

	return names
}

func Retrieve(commands []*Command, name string) *Command {
	for _, c := range commands {
		if slices.Contains(c.Aliases, name) {
			return c
		}
	}

	return nil
}
