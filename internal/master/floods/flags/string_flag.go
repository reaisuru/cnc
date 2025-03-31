package flags

import (
	"cnc/internal/clients/packet"
	"fmt"
	"strconv"
)

type StringFlag struct {
	maxLength int
}

// NewStringFlag makes a new string flag
func NewStringFlag(maxLength int) *StringFlag {
	return &StringFlag{maxLength: maxLength}
}

// Validate attempts to validate a string flag type.
func (s *StringFlag) Validate(literal string, parent *Flag) error {
	if len(literal) > s.maxLength {
		return fmt.Errorf("flag %s exceeds a max length of %d", strconv.Quote(parent.Name), s.maxLength)
	}

	return nil
}

func (s *StringFlag) Write(literal string, packet *packet.Packet, parent *Flag) error {
	packet.AddString(literal)
	return nil
}

func (s *StringFlag) Name() string {
	return "string"
}

func (s *StringFlag) TypeID() uint8 {
	return TypeString
}
