package crypto

import "testing"

func TestSigVerify(t *testing.T) {
	privKey := GenPrivKey()
	pubKey := privKey.PubKey()
	msg := RandBytes(257)
	sig := privKey.Sign(msg)
	if valid := pubKey.Verify(msg, sig); !valid {
		t.Error("invalid signature")
	}
}
