package flags

import (
	"cnc/internal/clients/packet"
	"math/rand"
	"strconv"
)

type BooleanFlag struct{}

// NewBooleanFlag makes a new boolean flag
func NewBooleanFlag() *BooleanFlag {
	return &BooleanFlag{}
}

// Validate attempts to validate a string flag type.
func (s *BooleanFlag) Validate(literal string, parent *Flag) error {
	return nil
}

func (s *BooleanFlag) Write(literal string, packet *packet.Packet, parent *Flag) error {
	packet.AddUint16(1)

	if literal == "maybe" {
		packet.AddUint8(uint8(rand.Int() % 2))
		return nil
	}

	yes, err := strconv.ParseBool(literal)
	if err != nil {
		return err
	}

	// trash
	if yes {
		packet.AddUint8(1)
	} else {
		packet.AddUint8(0)
	}

	return nil
}

func (s *BooleanFlag) Name() string {
	return "bool"
}

func (s *BooleanFlag) TypeID() uint8 {
	return FlagBool
}
