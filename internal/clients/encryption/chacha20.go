package encryption

import (
	"encoding/binary"
)

func u32t8le(v uint32, p []byte) {
	binary.LittleEndian.PutUint32(p, v)
}

func u8t32le(p []byte) uint32 {
	return binary.LittleEndian.Uint32(p)
}

func rotl32(x uint32, n int) uint32 {
	return (x << n) | (x >> (32 - n))
}

func chacha20Quarterround(x *[16]uint32, a, b, c, d int) {
	x[a] += x[b]
	x[d] = rotl32(x[d]^x[a], 16)
	x[c] += x[d]
	x[b] = rotl32(x[b]^x[c], 12)
	x[a] += x[b]
	x[d] = rotl32(x[d]^x[a], 8)
	x[c] += x[d]
	x[b] = rotl32(x[b]^x[c], 7)
}

func chacha20Serialize(in *[16]uint32, output []byte) {
	for i := 0; i < 16; i++ {
		u32t8le(in[i], output[i*4:(i+1)*4])
	}
}

func chacha20Block(in *[16]uint32, out []byte, numRounds int) {
	var x [16]uint32
	copy(x[:], in[:])

	for i := numRounds; i > 0; i -= 2 {
		chacha20Quarterround(&x, 0, 4, 8, 12)
		chacha20Quarterround(&x, 1, 5, 9, 13)
		chacha20Quarterround(&x, 2, 6, 10, 14)
		chacha20Quarterround(&x, 3, 7, 11, 15)
		chacha20Quarterround(&x, 0, 5, 10, 15)
		chacha20Quarterround(&x, 1, 6, 11, 12)
		chacha20Quarterround(&x, 2, 7, 8, 13)
		chacha20Quarterround(&x, 3, 4, 9, 14)
	}

	for i := 0; i < 16; i++ {
		x[i] += in[i]
	}

	chacha20Serialize(&x, out)
}

func chacha20InitState(s *[16]uint32, key []byte, counter uint32, nonce []byte) {
	// convert magic number to string: "expand 32-byte k"
	s[0] = 0x61707865
	s[1] = 0x3320646e
	s[2] = 0x79622d32
	s[3] = 0x6b206574

	for i := 0; i < 8; i++ {
		s[4+i] = u8t32le(key[i*4 : (i+1)*4])
	}

	s[12] = counter

	for i := 0; i < 3; i++ {
		s[13+i] = u8t32le(nonce[i*4 : (i+1)*4])
	}
}

func Chacha20Encrypt(key []byte, counter uint32, nonce []byte, in []byte) []byte {
	var out []byte
	out = append(out, in...)
	Chacha20(key, counter, nonce, in, out)
	return out
}

func Chacha20(key []byte, counter uint32, nonce []byte, in, out []byte) {
	var s [16]uint32
	var block [64]byte

	chacha20InitState(&s, key[:], counter, nonce[:])

	inLen := len(in)
	for i := 0; i < inLen; i += 64 {
		chacha20Block(&s, block[:], 20)
		s[12]++

		for j := i; j < i+64 && j < inLen; j++ {
			out[j] = in[j] ^ block[j-i]
		}
	}
}
