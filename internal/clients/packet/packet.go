package packet

import (
	"encoding/binary"
	"strings"
)

type Packet struct {
	data []byte
	size uint16
}

// New creates a new packet instance
func New() *Packet {
	return new(Packet)
}

// AddBoolean will add a boolean to the buffer.
func (p *Packet) AddBoolean(d bool) *Packet {
	if d {
		p.data = append(p.data, 1)
	} else {
		p.data = append(p.data, 0)
	}

	p.size += 1
	return p
}

// AddUint8 will add an unsigned 8-bit integer to the buffer.
func (p *Packet) AddUint8(d uint8) *Packet {
	p.data = append(p.data, d)
	p.size += 1
	return p
}

// AddUint16 will add an unsigned 16-bit integer to the buffer.
func (p *Packet) AddUint16(d uint16) *Packet {
	var buffer = make([]byte, 2)
	binary.BigEndian.PutUint16(buffer, d)

	p.data = append(p.data, buffer...)
	p.size += 2
	return p
}

// AddUint32 will add an unsigned 32-bit integer to the buffer.
func (p *Packet) AddUint32(d uint32) *Packet {
	var buffer = make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, d)

	p.data = append(p.data, buffer...)
	p.size += 4
	return p
}

// AddBytes will add bytes to the buffer.
func (p *Packet) AddBytes(d []byte) *Packet {
	var buffer = make([]byte, 2)

	binary.BigEndian.PutUint16(buffer, uint16(len(d)))
	p.data = append(p.data, buffer...)
	p.data = append(p.data, d...)

	p.size += 2
	p.size += uint16(len(d))

	return p
}

// AddString will add a string to the buffer.
func (p *Packet) AddString(d string) *Packet {
	return p.AddBytes(encodeUTF8(d))
}

// Size returns the size of a packet
func (p *Packet) Size() uint16 {
	return p.size
}

// Bytes returns the bytes of the buffer.
func (p *Packet) Bytes() []byte {
	return p.data
}

func encodeUTF8(s string) []byte {
	if !strings.HasSuffix(s, "\x00") {
		s += "\x00"
	}

	return []byte(s)
}
