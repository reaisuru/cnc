package command

import (
	"cnc/internal/database"
	"cnc/internal/master/floods"
	"fmt"
	"strings"
	"unicode/utf8"
)

type AutoCompleter struct {
	// fully entered string
	full string

	// spaceIndex commandIndex
	index map[int]int

	// index value
	commands map[int]string

	// arguments
	arguments map[int]int

	// the database profile of the user
	profile *database.UserProfile
}

func NewCompleter(profile *database.UserProfile) *AutoCompleter {
	return &AutoCompleter{
		index:     make(map[int]int),
		full:      "",
		commands:  make(map[int]string),
		arguments: make(map[int]int),
		profile:   profile,
	}
}

func (a *AutoCompleter) AutoComplete(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	if key != '\t' || pos != len(line) {
		return
	}

	parts := strings.Split(line, " ")
	partCount := len(parts)

	// root command
	if partCount == 1 {
		newLine = a.handleCommands(Commands, 0, line)
		newPos = utf8.RuneCountInString(newLine)
		ok = true
		return
	}

	// retrieve the base command
	baseCommand := Retrieve(Commands, parts[0])
	if baseCommand == nil {
		return
	}

	newLine = parts[0]
	depth := 1

	for i := 1; i < partCount; i++ {
		if i == partCount-1 {
			if len(baseCommand.SubCommands) > 0 {
				newLine += " " + a.handleCommands(baseCommand.SubCommands, i, parts[i])
			} else {
				newLine += " " + a.handleArguments(baseCommand, parts[i], i-depth)
			}
		} else {
			subCommand := Retrieve(baseCommand.SubCommands, parts[i])

			if subCommand != nil {
				newLine += " " + subCommand.Aliases[0]
				baseCommand = subCommand
				depth++
			} else {
				newLine += " " + parts[i]
			}
		}
	}

	// check if newline is not invalid if no suggestions are available
	if newLine == line {
		return
	}

	newPos = utf8.RuneCountInString(newLine)
	ok = true
	return
}

func (a *AutoCompleter) handleCommands(commands []*Command, spaceCount int, line string) (newLine string) {
	if len(commands) == 0 {
		return line
	}

	var rangeIndex = a.index[spaceCount]
	var possibleCommands = Names(commands, a.profile)

	if spaceCount <= 0 {
		for v, vec := range floods.VectorList {
			if !a.profile.ContainsRole(vec.Roles) {
				continue
			}

			possibleCommands = append(possibleCommands, v)
		}
	}

	var cmd string

	// rangeIndex > all possible commands, reset
	if rangeIndex >= len(possibleCommands) {
		rangeIndex = 0
	}

	// check if user typed out a part of the cmd ig
	if a.commands[spaceCount] != strings.TrimSpace(line) {
		closestCommand, foundIndex := closestString(possibleCommands, line)

		// if there is command that's close to the string entered, we set it to that
		if closestCommand != "" {
			cmd = closestCommand
			rangeIndex = foundIndex
		} else {
			cmd = line // else we'll just set the line to the other thing
		}

		a.arguments = make(map[int]int) // reset
	} else {
		// cycle through commands
		cmd = possibleCommands[rangeIndex]
		a.arguments = make(map[int]int) // reset
	}

	// save all stuff
	a.index[spaceCount] = (rangeIndex + 1) % len(possibleCommands)
	a.commands[spaceCount] = cmd

	return cmd
}

func (a *AutoCompleter) handleArguments(command *Command, line string, spaceCount int) (newLine string) {
	var cmd string

	if spaceCount <= 0 {
		spaceCount = 0
	}

	// no subcommands enabled
	if len(command.Arguments) <= 0 {
		return cmd
	}

	if spaceCount >= len(command.Arguments) {
		return cmd
	}

	arg := command.Arguments[spaceCount]
	if arg == nil {
		return cmd
	}

	switch arg.Type {
	case ArgumentBoolean:
		return fmt.Sprintf("%v", a.doArgIndex(spaceCount, 1) > 0)
	case ArgumentUser:
		users, err := database.User.SelectAll()
		if err != nil {
			return cmd
		}

		index := a.doArgIndex(spaceCount, len(users)-1)
		return users[index].Name
	case ArgumentInteger:
		return fmt.Sprintf("%v", a.doArgIndex(spaceCount, 0))
	default:
	}

	return cmd
}

func closestString(str []string, value string) (string, int) {
	for i, s := range str {
		if strings.HasPrefix(s, value) && len(value) >= 1 {
			return s, i
		}
	}

	return "", -1
}

// doArgIndex is literal trash code. but works perfectly fine
func (a *AutoCompleter) doArgIndex(spaceCount, max int) int {
	i, ok := a.arguments[spaceCount]

	if !ok || (i >= max && max != 0) {
		a.arguments[spaceCount] = 0
		i = 0
	} else {
		i++
		a.arguments[spaceCount] = i
	}

	return i
}
