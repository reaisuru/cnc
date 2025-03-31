package flags

import (
	"cnc/internal/clients/packet"
	"encoding/hex"
	"fmt"
	"strconv"
)

type HexFlag struct {
	maxLength int
}

// NewHexFlag makes a new string flag
func NewHexFlag(maxLength int) *HexFlag {
	return &HexFlag{maxLength: maxLength}
}

// Validate attempts to validate a string flag type.
func (s *HexFlag) Validate(literal string, parent *Flag) error {
	bytes, err := hex.DecodeString(literal)
	if err != nil {
		return fmt.Errorf("flag %s is an invalid hexadecimal stream", strconv.Quote(parent.Name))
	}

	if len(bytes) > s.maxLength {
		return fmt.Errorf("flag %s exceeds a max length of %d", strconv.Quote(parent.Name), s.maxLength)
	}

	return nil
}

func (s *HexFlag) Write(literal string, packet *packet.Packet, parent *Flag) error {
	bytes, _ := hex.DecodeString(literal)
	packet.AddBytes(bytes)
	return nil
}

func (s *HexFlag) Name() string {
	return "hex"
}

func (s *HexFlag) TypeID() uint8 {
	return TypeHex
}
