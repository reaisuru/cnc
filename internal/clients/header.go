package clients

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Header struct {
	// Op is the instruction ID
	Op uint8

	// Padding
	_ uint8

	// Len is the body length
	Len uint16
}

// Bytes converts the header into bytes.
func (h *Header) Bytes() []byte {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, h)
	return buf.Bytes()
}

// String converts the header into a printable format
func (h *Header) String() string {
	return fmt.Sprintf("Header{op=%d, len=%d}", h.Op, h.Len)
}
