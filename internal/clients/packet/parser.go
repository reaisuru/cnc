package packet

import (
	"bytes"
	"encoding/binary"
)

type Parser struct {
	buffer []byte
}

// NewParser creates a new Parser with the provided buffer.
func NewParser(buffer []byte) *Parser {
	return &Parser{buffer: buffer}
}

// ParseInt8 parses the next byte as an int8.
func (p *Parser) ParseInt8() int8 {
	if p.Length() < 1 {
		return 0
	}

	value := p.buffer[:1][0]
	p.buffer = p.buffer[1:]

	return int8(value)
}

// ParseBool parses a boolean
func (p *Parser) ParseBool() bool {
	return p.ParseInt8() == 1
}

// ParseInt16 parses the next 2 bytes as an int16
func (p *Parser) ParseInt16() int16 {
	if p.Length() < 2 {
		return 0
	}

	value := binary.BigEndian.Uint16(p.buffer[:2])
	p.buffer = p.buffer[2:]

	return int16(value)
}

// ParseInt32 parses the next 4 bytes as an int32
func (p *Parser) ParseInt32() int {
	if p.Length() < 4 {
		return 0
	}

	value := int(binary.BigEndian.Uint32(p.buffer[:4]))
	p.buffer = p.buffer[4:]

	return value
}

// ParseInt64 parses the next 8 bytes as an int64
func (p *Parser) ParseInt64() int64 {
	if p.Length() < 8 {
		return 0
	}

	value := int64(binary.BigEndian.Uint64(p.buffer[:8]))
	p.buffer = p.buffer[8:]

	return value
}

// ParseBytes parses the next 4 bytes to determine the length of the byte slice,
// then reads that many bytes from the buffer.
func (p *Parser) ParseBytes() []byte {
	if p.Length() < 4 {
		return nil
	}

	byteSize := p.ParseInt32()
	if byteSize > p.Length() {
		byteSize = p.Length()
	}

	bytesBuffer := p.buffer[:byteSize]
	p.buffer = p.buffer[byteSize:]

	return bytesBuffer
}

// ParseAtLeastBytes parses a specified number of bytes
func (p *Parser) ParseAtLeastBytes(numBytes int) []byte {
	if numBytes > p.Length() {
		numBytes = p.Length()
	}

	bytesBuffer := p.buffer[:numBytes]
	p.buffer = p.buffer[numBytes:]

	return bytesBuffer
}

// ParseString parses bytes into a string and removes the null terminator
func (p *Parser) ParseString() string {
	return StripNull(string(p.ParseBytes()))
}

// Length returns the number of unread bytes
func (p *Parser) Length() int {
	return len(p.buffer)
}

// Buffer returns the unread bytes
func (p *Parser) Buffer() []byte {
	return p.buffer
}

// Readable checks if the buffer has enough bytes to read the given types
func (p *Parser) Readable(types ...any) bool {
	bytesRead := 0
	totalSize := p.Length()

	for _, t := range types {
		var size int

		switch t.(type) {
		case *byte, *bool:
			size = 1
		case *uint16:
			size = 2
		case *uint32:
			size = 4
		case *uint64:
			size = 8
		case *[]byte, *string:
			if totalSize-bytesRead < 4 {
				return false
			}

			length := int(binary.BigEndian.Uint32(p.buffer[bytesRead : bytesRead+4]))
			bytesRead += 4

			if totalSize-bytesRead < length {
				return false
			}

			size = length
		default:
			return false
		}

		if totalSize-bytesRead < size {
			return false
		}

		bytesRead += size
	}

	return true
}

func StripNull(s string) string {
	return string(bytes.Trim([]byte(s), "\x00"))
}
