package crypto

import (
	"bytes"
	"testing"
)

func TestSigVerify(t *testing.T) {
	privKey := GenPrivKey()
	pubKey := privKey.PubKey()
	msg := RandBytes(64)
	sig := privKey.Sign(msg)
	if valid := pubKey.Verify(msg, sig); !valid {
		t.Error("invalid signature")
	}
}

func TestPubKey(t *testing.T) {
	privKey := GenPrivKey()
	newPrivKeyP := new(PrivateKey)
	copy(newPrivKeyP[:], privKey[:32])
	newPrivKeyP.SetPubKey()
	if bytes.Compare(privKey[32:], newPrivKeyP[32:]) != 0 {
		t.Error("new PublicKey doesn't match")
	}
}
