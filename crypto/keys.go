package crypto

import (
	"fmt"

	"golang.org/x/crypto/ed25519"
)

type PublicKey [32]byte
type PrivateKey [64]byte

type Signature [64]byte

func GenPrivKey() PrivateKey {
	privKey := PrivateKey{}
	edPrivKey := ed25519.NewKeyFromSeed(RandBytes(32))
	copy(privKey[:], edPrivKey)
	return privKey
}

func (privKey PrivateKey) PubKey() PublicKey {
	pubKey := PublicKey{}
	privKeyBytes := privKey[:]

	for _, v := range privKeyBytes[32:] {
		if v != 0 {
			copy(pubKey[:], privKeyBytes[32:])
			return pubKey
		}
	}

	pubKeyBytes := ed25519.PrivateKey(privKeyBytes).Public().([]byte)
	copy(pubKey[:], pubKeyBytes)
	return pubKey
}

func (privKey PrivateKey) Sign(msg []byte) Signature {
	sig := Signature{}
	edPrivKey := ed25519.PrivateKey(privKey[:])
	sigBytes := ed25519.Sign(edPrivKey, msg)
	copy(sig[:], sigBytes)
	return sig
}

func (pubKey PublicKey) Verify(msg []byte, sig Signature) bool {
	edPubKey := ed25519.PublicKey(pubKey[:])
	return ed25519.Verify(edPubKey, msg, sig[:])
}

func (pubKey PublicKey) String() string {
	return fmt.Sprintf("%x", pubKey[:])
}

func (pubKey PublicKey) StringWithDashes() string {
	x := fmt.Sprintf("%x", pubKey[:])
	return fmt.Sprintf("%v-%v-%v-%v", x[:16], x[16:32], x[32:48], x[48:])
}

func (pubKey PublicKey) StringCompact() string {
	x := fmt.Sprintf("%x", pubKey[:])
	return fmt.Sprintf("%v...%v", x[:4], x[60:])
}
