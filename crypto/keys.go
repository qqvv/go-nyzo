package crypto

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/ed25519"
)

type PublicKey [32]byte
type PrivateKey [64]byte // PrivateKey + PublicKey

type Signature [64]byte

func GenPrivKey() PrivateKey {
	privKey := PrivateKey{}
	edPrivKey := ed25519.NewKeyFromSeed(RandBytes(32))
	copy(privKey[:], edPrivKey)
	return privKey
}

// Set the PublicKey in PrivateKey if missing
func (privKey *PrivateKey) SetPubKey() {
	privKeyBytes := privKey[:]

	for _, v := range privKeyBytes[32:] {
		if v != 0 {
			return
		}
	}

	edPrivKey := ed25519.NewKeyFromSeed(privKey[:32])
	copy(privKey[:], edPrivKey[:])
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

	edPrivKey := ed25519.NewKeyFromSeed(privKey[:32])
	copy(pubKey[:], edPrivKey[32:])
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

func (privKey PrivateKey) String() string {
	return fmt.Sprintf("%x", privKey[:32])
}

func (sig Signature) String() string {
	return fmt.Sprintf("%x", sig[:])
}

func (pubKey PublicKey) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(pubKey.String())
	return b, err
}

func (privKey PrivateKey) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(privKey.String())
	return b, err
}

func (sig Signature) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(sig.String())
	return b, err
}

func (pubKey PublicKey) UnmarshalJSON(data []byte) error {
	var x string
	err := json.Unmarshal(data, x)
	if err != nil {
		return err
	}
	b, err := hex.DecodeString(x)
	if err != nil {
		return err
	}
	copy(pubKey[:], b)
	return nil
}

func (privKey PrivateKey) UnmarshalJSON(data []byte) error {
	var x string
	err := json.Unmarshal(data, x)
	if err != nil {
		return err
	}
	b, err := hex.DecodeString(x)
	if err != nil {
		return err
	}
	copy(privKey[:], b)
	return nil
}

func (sig Signature) UnmarshalJSON(data []byte) error {
	var x string
	err := json.Unmarshal(data, x)
	if err != nil {
		return err
	}
	b, err := hex.DecodeString(x)
	if err != nil {
		return err
	}
	copy(sig[:], b)
	return nil
}
