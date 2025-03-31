package encryption

import "crypto/sha256"

const (
	HmacBlockSize = 64
)

// HmacSha256 will generate a hmac checksum based on sha256
func HmacSha256(key []byte, data []byte) []byte {
	var keyBlock [HmacBlockSize]byte
	var ipad [HmacBlockSize]byte
	var opad [HmacBlockSize]byte

	if len(key) > HmacBlockSize {
		hash := sha256.Sum256(key)
		copy(keyBlock[:], hash[:])
	} else {
		copy(keyBlock[:], key)
	}

	for i := 0; i < HmacBlockSize; i++ {
		ipad[i] = keyBlock[i] ^ 0x36
		opad[i] = keyBlock[i] ^ 0x5c
	}

	innerHash := sha256.New()
	innerHash.Write(ipad[:])
	innerHash.Write(data)

	temp := innerHash.Sum(nil)

	outerHash := sha256.New()
	outerHash.Write(opad[:])
	outerHash.Write(temp)

	return outerHash.Sum(nil)
}
