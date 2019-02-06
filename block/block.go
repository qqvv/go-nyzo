package block

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/qqvv/go-nyzo/balancelist"
	"github.com/qqvv/go-nyzo/crypto"
	"github.com/qqvv/go-nyzo/transaction"

	"github.com/davecgh/go-spew/spew"
)

type Block struct {
	Height                int64             `json:"height"`
	PrevBlockHash         crypto.Hash       `json:"prevBlockHash"`
	StartTimestamp        int64             `json:"startTimestamp"`
	VerificationTimestamp int64             `json:"verificationTimestamp"`
	Transactions          []*transaction.Tx `json:"transactions"`
	BalancelistHash       crypto.Hash       `json:"balancelistHash"`
	VerifierID            crypto.PublicKey  `json:"verifierId"`
	VerifierSig           crypto.Signature  `json:"verifierSig"`
}

func (bl *Block) Serialize() []byte {
	buf := bytes.NewBuffer(bl.ForSigning())
	binary.Write(buf, binary.BigEndian, bl.VerifierSig)

	return buf.Bytes()
}

func (bl *Block) ForSigning() []byte {
	buf := new(bytes.Buffer)
	buf.Grow(bl.Size())

	binary.Write(buf, binary.BigEndian, bl.Height)
	binary.Write(buf, binary.BigEndian, bl.PrevBlockHash)
	binary.Write(buf, binary.BigEndian, bl.StartTimestamp/1000/1000)        // to milli
	binary.Write(buf, binary.BigEndian, bl.VerificationTimestamp/1000/1000) // to milli

	binary.Write(buf, binary.BigEndian, int32(len(bl.Transactions)))
	for _, tx := range bl.Transactions {
		binary.Write(buf, binary.BigEndian, tx.Serialize())
	}

	binary.Write(buf, binary.BigEndian, bl.BalancelistHash)
	binary.Write(buf, binary.BigEndian, bl.VerifierID)

	return buf.Bytes()
}

func (bl *Block) Deserialize(i interface{}) error {
	var buf *bytes.Buffer

	switch i.(type) {
	case *bytes.Buffer:
		buf = i.(*bytes.Buffer)
	case []byte:
		buf = bytes.NewBuffer(i.([]byte))
	default:
		return fmt.Errorf("cannot deserialize block from %#v", i)
	}

	blockCount := int16(0)
	binary.Read(buf, binary.BigEndian, &blockCount)

	// TODO: handle multiple blocks
	if blockCount != 1 {
		return fmt.Errorf("cannot deserialize %v blocks", blockCount)
	}

	binary.Read(buf, binary.BigEndian, &bl.Height)
	binary.Read(buf, binary.BigEndian, &bl.PrevBlockHash)
	binary.Read(buf, binary.BigEndian, &bl.StartTimestamp)
	binary.Read(buf, binary.BigEndian, &bl.VerificationTimestamp)
	bl.StartTimestamp *= 1000 * 1000        // to nano
	bl.VerificationTimestamp *= 1000 * 1000 // to nano

	txCount := int(binary.BigEndian.Uint32(buf.Next(4)))
	for i := 0; i < txCount; i++ {
		tx := &transaction.Tx{}
		if err := tx.Deserialize(buf); err != nil {
			return fmt.Errorf("error deserializing Tx from block: %v", err)
		}
		bl.Transactions = append(bl.Transactions, tx)
	}

	binary.Read(buf, binary.BigEndian, &bl.BalancelistHash)
	binary.Read(buf, binary.BigEndian, &bl.VerifierID)
	binary.Read(buf, binary.BigEndian, &bl.VerifierSig)

	balList := &balancelist.List{}
	balList.Deserialize(buf)
	spew.Dump(balList)

	return nil
}

func (bl *Block) Size() int {
	size := 124
	for _, tx := range bl.Transactions {
		size += tx.SerializedLen()
	}
	if !bytes.Equal(bl.VerifierSig[:], make([]byte, 64)) {
		size += 64
	}
	return size
}

func (bl *Block) Sign(privKey crypto.PrivateKey) {
	bl.VerifierID = privKey.PubKey()
	bl.VerifierSig = privKey.Sign(bl.ForSigning())
}

func (bl *Block) Hash() crypto.Hash {
	return crypto.DoubleSHA256(bl.VerifierSig[:])
}

func (bl *Block) PrevHeight() int64 {
	return bl.Height - 1
}

func New(height, startTimestamp int64, prevBlockHash, balancelistHash crypto.Hash) *Block {
	bl := &Block{}
	bl.Height = height
	bl.PrevBlockHash = prevBlockHash
	bl.StartTimestamp = startTimestamp
	bl.BalancelistHash = balancelistHash
	return bl
}
