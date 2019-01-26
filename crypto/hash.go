package crypto

import (
	"crypto/sha256"
	"fmt"
)

type Hash [32]byte

func (h *Hash) String() string {
	return fmt.Sprintf("%X", h[:])
}

func DoubleSHA256(data []byte) Hash {
	h := sha256.Sum256(data)
	h = sha256.Sum256(h[:])
	return Hash(h)
}
