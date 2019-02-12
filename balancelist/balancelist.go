package balancelist

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/qqvv/go-nyzo/crypto"
)

type List struct {
	Height        int64              `json:"height"`
	RolloverFees  byte               `json:"rolloverFees"`
	PrevVerifiers []crypto.PublicKey `json:"prevVerifiers"`
	Items         []*Item            `json:"items"`
}

type Item struct {
	ID             crypto.PublicKey `json:"id"`
	Balance        int64            `json:"balance"`
	BlocksUntilFee int16            `json:"blocksUntilFee"`
}

func (list *List) Serialize() []byte {
	buf := new(bytes.Buffer)
	buf.Grow(list.SerializedLen())

	binary.Write(buf, binary.BigEndian, list.Height)
	binary.Write(buf, binary.BigEndian, list.RolloverFees)

	for _, id := range list.PrevVerifiers {
		binary.Write(buf, binary.BigEndian, id)
	}

	binary.Write(buf, binary.BigEndian, int32(len(list.Items)))
	for _, item := range list.Items {
		binary.Write(buf, binary.BigEndian, item.ID)
		binary.Write(buf, binary.BigEndian, item.Balance)
		binary.Write(buf, binary.BigEndian, item.BlocksUntilFee)
	}

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

func (list *List) SerializedLen() int {
	size := 13                           // height (8) + rolloverfees (1) + length (4)
	size += len(list.PrevVerifiers) * 32 // pubk (32)
	size += len(list.Items) * 42         // pubk (32) + balance (8) + blocksuntilfee (2)
	return size
}

func NewList() *List {
	return &List{}
}
