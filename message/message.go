package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/qqvv/go-nyzo/crypto"
)

type Msg struct {
	Timestamp int64
	Type      MsgType
	Content   []byte
	ID        crypto.PublicKey
	Sig       crypto.Signature

	sourceIP []byte
}

func (msg *Msg) Serialize() []byte {
	buf := new(bytes.Buffer)
	buf.Grow(msg.SerializedLen())

	binary.Write(buf, binary.BigEndian, int32(msg.SerializedLen()))
	binary.Write(buf, binary.BigEndian, msg.ForSigning())
	binary.Write(buf, binary.BigEndian, msg.Sig)

	return buf.Bytes()
}

func (msg *Msg) ForSigning() []byte {
	buf := new(bytes.Buffer)
	buf.Grow(msg.SerializedLen() - 4) // without length

	binary.Write(buf, binary.BigEndian, msg.Timestamp/1000/1000) // to milli
	binary.Write(buf, binary.BigEndian, msg.Type)
	binary.Write(buf, binary.BigEndian, msg.Content)
	binary.Write(buf, binary.BigEndian, msg.ID)

	return buf.Bytes()
}

func (msg *Msg) Deserialize(i interface{}) error {
	var buf *bytes.Buffer

	switch i.(type) {
	case *bytes.Buffer:
		buf = i.(*bytes.Buffer)
	case []byte:
		buf = bytes.NewBuffer(i.([]byte))
	default:
		return fmt.Errorf("cannot deserialize Msg from %#v", i)
	}

	binary.Read(buf, binary.BigEndian, &msg.Timestamp)
	msg.Timestamp *= 1000 * 1000 // to nano
	binary.Read(buf, binary.BigEndian, &msg.Type)
	binary.Read(buf, binary.BigEndian, &msg.Content)
	binary.Read(buf, binary.BigEndian, &msg.ID)
	binary.Read(buf, binary.BigEndian, &msg.Sig)

	return nil
}

func (msg *Msg) Sign(privKey crypto.PrivateKey) {
	msg.ID = privKey.PubKey()
	msg.Sig = privKey.Sign(msg.ForSigning())
}

func (msg *Msg) SerializedLen() int {
	return 110 + len(msg.Content)
}

func (msg *Msg) VerifySig() bool {
	return msg.ID.Verify(msg.ForSigning(), msg.Sig)
}

func New(msgType int) *Msg {
	msg := &Msg{
		Timestamp: time.Now().UnixNano(),
		Type:      MsgType(msgType),
	}
	return msg
}
