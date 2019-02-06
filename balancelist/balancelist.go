package balancelist

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/qqvv/go-nyzo/crypto"
)

type List struct {
	Height        int64
	RolloverFees  byte
	PrevVerifiers []crypto.PublicKey
	Items         []*Item
}

type Item struct {
	ID             crypto.PublicKey
	Balance        int64
	BlocksUntilFee int16
}

func (list *List) Serialize() []byte {
	buf := new(bytes.Buffer)
	// TODO
	return buf.Bytes()
}

func (list *List) Deserialize(i interface{}) error {
	var buf *bytes.Buffer

	switch i.(type) {
	case *bytes.Buffer:
		buf = i.(*bytes.Buffer)
	case []byte:
		buf = bytes.NewBuffer(i.([]byte))
	default:
		return fmt.Errorf("cannot deserialize balancelist from %#v", i)
	}

	binary.Read(buf, binary.BigEndian, &list.Height)
	binary.Read(buf, binary.BigEndian, &list.RolloverFees)

	prevVerifierCount := 9
	if list.Height < int64(9) {
		prevVerifierCount = int(list.Height)
	}
	for i := 0; i < prevVerifierCount; i++ {
		id := crypto.PublicKey{}
		binary.Read(buf, binary.BigEndian, &id)
		list.PrevVerifiers = append(list.PrevVerifiers, id)
	}

	var itemCount int32
	binary.Read(buf, binary.BigEndian, &itemCount)
	for i := 0; int32(i) < itemCount; i++ {
		item := &Item{}
		binary.Read(buf, binary.BigEndian, &item.ID)
		binary.Read(buf, binary.BigEndian, &item.Balance)
		binary.Read(buf, binary.BigEndian, &item.BlocksUntilFee)
		list.Items = append(list.Items, item)
	}

	return nil
}

func NewList() *List {
	return &List{}
}
