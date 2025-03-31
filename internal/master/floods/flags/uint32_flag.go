package flags

import (
	"cnc/internal/clients/packet"
	"fmt"
	"strconv"
)

type Uint32Flag struct {
	minimum int
	maximum int
}

// NewUint32Flag makes a new uint32 flag
func NewUint32Flag(minimum, maximum int) *Uint32Flag {
	return &Uint32Flag{
		minimum: minimum,
		maximum: maximum,
	}
}

// Validate attempts to validate a string flag type.
func (s *Uint32Flag) Validate(literal string, parent *Flag) error {
	integer, err := strconv.Atoi(literal)
	if err != nil {
		return fmt.Errorf("invalid integer input for %s: %s", strconv.Quote(parent.Name), strconv.Quote(literal))
	}

	if integer < 0 || integer > 2_147_483_647 {
		return fmt.Errorf("%s must be within the 16-bit unsigned integer range (0 - 2147483647)", parent.Name)
	}

	if s.minimum != -1 && integer < s.minimum {
		return fmt.Errorf("%s must be within the minimum %d", parent.Name, s.minimum)
	} else if s.maximum != -1 && integer > s.maximum {
		return fmt.Errorf("%s must be within the maximum %d", parent.Name, s.minimum)
	}

	return nil
}

func (s *Uint32Flag) Write(literal string, packet *packet.Packet, parent *Flag) error {
	integer, _ := strconv.Atoi(literal)
	packet.AddUint16(4) // uint32 = 4 bytes
	packet.AddUint32(uint32(integer))
	return nil
}

func (s *Uint32Flag) Name() string {
	return "uint32"
}

func (s *Uint32Flag) TypeID() uint8 {
	return FlagUint32
}
