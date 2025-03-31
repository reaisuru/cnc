package flags

import (
	"cnc/internal/clients/packet"
)

const (
	TypeString uint8 = iota
	TypeHex
	FlagBool
	FlagUint8
	FlagUint16
	FlagUint32
)

type Flag struct {
	// ID is the id of the current flag
	ID uint8

	// Name is the name of the flag
	Name string

	// Description is a simple and short description for an attack flag
	Description string

	// Type is the flag type.
	Type FlagType

	// Options are the options
	Options *FlagOptions
}

type FlagOptions struct {
	Clientside bool
	Admin      bool
	Invisible  bool
}

type FlagOption func(*FlagOptions) error

type FlagType interface {
	Validate(literal string, parent *Flag) error
	Write(literal string, packet *packet.Packet, parent *Flag) error
	Name() string
	TypeID() uint8
}
