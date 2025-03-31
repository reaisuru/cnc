package command

import (
	"fmt"
	"strings"
)

// NewContext creates a new command context with parsed arguments
func NewContext(parent *Command, cmd *Command, envs []string, arguments ...string) (*Context, error) {
	// create new command context
	var ctx = new(Context)
	ctx.arguments = make(map[string]*ParsedArgument)
	ctx.rawArgs = arguments
	ctx.parent = parent
	ctx.environments = make(map[string]string)

	// we'll just go through the envs I guess
	for _, env := range envs {
		prefix, suffix, ok := strings.Cut(env, "=")
		if !ok {
			continue
		}

		ctx.environments[prefix] = suffix
	}

	if cmd.Arguments == nil {
		return ctx, nil
	}

	// iterate through all registered arguments
	for pos, argument := range cmd.Arguments {
		// check if there are enough arguments provided
		if len(arguments) <= pos || len(arguments[pos:]) <= 0 {
			// if the argument is not required, check if there is a default one
			// and if there is not we just continue without doing anything...
			if !argument.Required {
				// check if we have a default value... we talked about that above.
				if argument.Default != nil {
					ctx.arguments[argument.Name] = &ParsedArgument{
						Type:  argument.Type,
						Value: argument.Default,
					}
				}

				continue
			}

			return nil, fmt.Errorf("missing required argument: %s", argument.Name)
		}

		// skip the iteration if the pos index is above the available arguments
		if pos >= len(arguments) {
			continue
		}

		// convert a raw string argument to a go type
		converted, err := argument.Literal2Go(arguments[pos:][0])
		if err != nil || converted == nil {
			return nil, err
		}

		// add parsed argument to context
		ctx.arguments[argument.Name] = &ParsedArgument{
			Type:  argument.Type,
			Value: converted,
		}
	}

	// return the context
	return ctx, nil
}
