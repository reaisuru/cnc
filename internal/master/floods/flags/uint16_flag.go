package flags

import (
	"cnc/internal/clients/packet"
	"fmt"
	"strconv"
)

type Uint16Flag struct {
	minimum int
	maximum int
}

// NewUint16Flag makes a new uint16 flag
func NewUint16Flag(minimum, maximum int) *Uint16Flag {
	return &Uint16Flag{
		minimum: minimum,
		maximum: maximum,
	}
}

// Validate attempts to validate a string flag type.
func (s *Uint16Flag) Validate(literal string, parent *Flag) error {
	integer, err := strconv.Atoi(literal)
	if err != nil {
		return fmt.Errorf("invalid integer input for %s: %s", strconv.Quote(parent.Name), strconv.Quote(literal))
	}

	if integer < 0 || integer > 65535 {
		return fmt.Errorf("%s must be within the 16-bit unsigned integer range (0 - 65535)", parent.Name)
	}

	if s.minimum != -1 && integer < s.minimum {
		return fmt.Errorf("%s must be within the minimum %d", parent.Name, s.minimum)
	} else if s.maximum != -1 && integer > s.maximum {
		return fmt.Errorf("%s must be within the maximum %d", parent.Name, s.minimum)
	}

	return nil
}

func (s *Uint16Flag) Write(literal string, packet *packet.Packet, parent *Flag) error {
	integer, _ := strconv.Atoi(literal)

	packet.AddUint16(2) // uint16 = 2 bytes
	packet.AddUint16(uint16(integer))

	return nil
}

func (s *Uint16Flag) Name() string {
	return "uint16"
}

func (s *Uint16Flag) TypeID() uint8 {
	return FlagUint16
}
