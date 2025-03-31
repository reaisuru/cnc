package sessions

import (
	"cnc/pkg/format"
	"cnc/pkg/sshd/terminal"
	"strconv"
	"time"
)

func (s *Session) Boolean(Prompt string, Default bool) (result bool, err error) {
	term := terminal.NewTerminal(s.Channel, Prompt)

	line, err := term.ReadLine()
	if err != nil {
		return false, err
	}

	if len(line) < 1 {
		return Default, nil
	}

	if line == "y" {
		return true, nil
	}

	if line == "n" {
		return false, nil
	}

	return strconv.ParseBool(line)
}

func (s *Session) Integer(Prompt string, Default int) (result int, err error) {
	term := terminal.NewTerminal(s.Channel, Prompt)

	line, err := term.ReadLine()
	if err != nil {
		return 0, err
	}

	if len(line) < 1 {
		return Default, nil
	}

	return strconv.Atoi(line)
}

func (s *Session) Time(Prompt string, Default string) (result time.Time, err error) {
	term := terminal.NewTerminal(s.Channel, Prompt)

	line, err := term.ReadLine()
	if err != nil {
		return time.Now(), err
	}

	if len(line) < 1 {
		duration, err := format.ModifiedParseDuration(Default)
		if err != nil {
			return time.Now(), err
		}

		return time.Now().Add(duration), nil
	}

	duration, err := format.ModifiedParseDuration(line)
	if err != nil {
		return time.Time{}, err
	}

	return time.Now().Add(duration), err
}

func (s *Session) Literal(Prompt string, Default string) (result string, err error) {
	term := terminal.NewTerminal(s.Channel, Prompt)

	line, err := term.ReadLine()
	if err != nil {
		return "", err
	}

	if len(line) < 1 {
		return Default, nil
	}

	return line, nil
}
