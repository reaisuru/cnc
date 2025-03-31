package flags

import (
	"cnc/internal/clients/packet"
	"fmt"
	"strconv"
)

type Uint8Flag struct{}

// NewUint8Flag makes a new uint8 flag
func NewUint8Flag() *Uint8Flag {
	return &Uint8Flag{}
}

// Validate attempts to validate a string flag type.
func (s *Uint8Flag) Validate(literal string, parent *Flag) error {
	integer, err := strconv.Atoi(literal)
	if err != nil {
		return fmt.Errorf("invalid integer input for %s: %s", strconv.Quote(parent.Name), strconv.Quote(literal))
	}

	if integer < 0 || integer > 255 {
		return fmt.Errorf("%s must be within the 8-bit unsigned integer range (0 - 255)", parent.Name)
	}

	return nil
}

func (s *Uint8Flag) Write(literal string, packet *packet.Packet, parent *Flag) error {
	integer, _ := strconv.Atoi(literal)
	packet.AddUint16(1) // uint8 = 1 byte
	packet.AddUint8(uint8(integer))
	return nil
}

func (s *Uint8Flag) TypeID() uint8 {
	return FlagUint8
}

func (s *Uint8Flag) Name() string {
	return "uint8"
}
