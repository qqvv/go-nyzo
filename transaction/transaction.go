package transaction

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/qqvv/go-nyzo/crypto"
)

type Tx struct {
	Type           byte // 0=coingeneration, 1=seed, 2=standard
	Timestamp      int64
	Amount         int64
	RecipientID    crypto.PublicKey
	PrevHashHeight int64
	PrevHash       crypto.Hash
	SenderID       crypto.PublicKey
	SenderData     []byte
	SenderSig      crypto.Signature
}

func (tx *Tx) Validate() (bool, error) {
	// A transaction is valid if:
	// (1) the type is correct
	// (1) the signature is correct
	// (2) the transaction amount is at least µ1
	// (3) the wallet has enough coins to send it (not checked here!)
	// (4) the previous-block hash is correct
	// (5) the sender and receiver are different
	// (6) the block for the specified timestamp is still open for processing

	var valid bool
	var err error

	if tx.Type != 1 && tx.Type != 2 {
		valid = false
		err = fmt.Errorf("only seed (1) and standard (2) txes are valid after block 0")
	}

	if valid && !tx.SenderID.Verify(tx.ForSigning(), tx.SenderSig) {
		valid = false
		err = fmt.Errorf("signature is not valid")
	}

	if valid && tx.Amount < 1 {
		valid = false
		err = fmt.Errorf("tx amount must be at least µ1")
	}

	if valid && tx.Type == 1 {
		if !bytes.Equal(tx.SenderID[:], tx.RecipientID[:]) {
			valid = false
			err = fmt.Errorf("sender and recipient must match for a seed tx")
		}
	}

	if valid && tx.Type == 2 {
		if bytes.Equal(tx.SenderID[:], tx.RecipientID[:]) {
			valid = false
			err = fmt.Errorf("sender and recipient must be different for a standard tx")
		}
	}

	return valid, err
}

func (tx *Tx) Serialize() []byte {
	buf := new(bytes.Buffer)
	buf.Grow(tx.SerializedLen())

	binary.Write(buf, binary.BigEndian, tx.Type)
	binary.Write(buf, binary.BigEndian, tx.Timestamp/1000/1000) // to milli
	binary.Write(buf, binary.BigEndian, tx.Amount)
	binary.Write(buf, binary.BigEndian, tx.RecipientID)
	binary.Write(buf, binary.BigEndian, tx.PrevHashHeight)
	binary.Write(buf, binary.BigEndian, tx.SenderID)
	binary.Write(buf, binary.BigEndian, byte(len(tx.SenderData)))
	binary.Write(buf, binary.BigEndian, tx.SenderData)
	binary.Write(buf, binary.BigEndian, tx.SenderSig)

	return buf.Bytes()
}

func (tx *Tx) ForSigning() []byte {
	buf := new(bytes.Buffer)
	// type (1) + timestamp (8) + amount (8) + recipient pubk (32) + prevhashheight (8)
	// + prevhash (32) + sender pubk (32) + hash of senderdata (32)
	buf.Grow(153)

	binary.Write(buf, binary.BigEndian, tx.Type)
	binary.Write(buf, binary.BigEndian, tx.Timestamp/1000/1000) // to milli
	binary.Write(buf, binary.BigEndian, tx.Amount)
	binary.Write(buf, binary.BigEndian, tx.RecipientID)
	binary.Write(buf, binary.BigEndian, tx.PrevHash)
	binary.Write(buf, binary.BigEndian, tx.SenderID)
	binary.Write(buf, binary.BigEndian, crypto.DoubleSHA256(tx.SenderData))

	return buf.Bytes()
}

func (tx *Tx) Deserialize(i interface{}) error {
	var buf *bytes.Buffer

	switch i.(type) {
	case *bytes.Buffer:
		buf = i.(*bytes.Buffer)
	case []byte:
		buf = bytes.NewBuffer(i.([]byte))
	default:
		return fmt.Errorf("cannot deserialize tx from %#v", i)
	}

	binary.Read(buf, binary.BigEndian, &tx.Type)

	// 0=coingeneration, 1=seed, 2=standard
	if tx.GetType() > 2 || tx.GetType() < 0 {
		return fmt.Errorf("unknown tx type: %v", int(tx.Type))
	}

	binary.Read(buf, binary.BigEndian, &tx.Timestamp)
	tx.Timestamp *= 1000 * 1000 // to nano
	binary.Read(buf, binary.BigEndian, &tx.Amount)
	binary.Read(buf, binary.BigEndian, &tx.RecipientID)

	// Coingeneration transactions don't have the last fields
	if int(tx.Type) == 0 {
		return nil
	}

	binary.Read(buf, binary.BigEndian, &tx.PrevHashHeight)
	// TODO: get hash for prevheight

	binary.Read(buf, binary.BigEndian, &tx.SenderID)

	dataLen, _ := buf.ReadByte()
	if int(dataLen) > 32 {
		// You might think to return here as a length over 32 should be
		// invalid but this is what the original implementation does
		dataLen = 32
	}
	tx.SenderData = make([]byte, int(dataLen))
	binary.Read(buf, binary.BigEndian, &tx.SenderData)

	binary.Read(buf, binary.BigEndian, &tx.SenderSig)

	return nil
}

func (tx *Tx) SerializedLen() int {
	// type (1) + timestamp (8) + amount (8) + recipient pubk (32)
	// + prevhashheight (8) + prevhash (32) + sender pubk (32)
	// + sender sig (64) + senderdata len (1) + senderdata (0-32)
	return 186 + len(tx.SenderData)
}

func (tx *Tx) Sign(privKey crypto.PrivateKey) {
	tx.SenderID = privKey.PubKey()
	tx.SenderSig = privKey.Sign(tx.ForSigning())
}

func (tx *Tx) GetType() int {
	return int(tx.Type)
}

func (tx *Tx) Hash() crypto.Hash {
	return crypto.DoubleSHA256(tx.Serialize())
}

func NewStandard(amount int64, recipientID crypto.PublicKey, data []byte) *Tx {
	tx := &Tx{
		Type:           byte(2),
		Timestamp:      time.Now().Unix(),
		Amount:         amount,
		RecipientID:    recipientID,
		PrevHashHeight: 0,
		SenderData:     data,
	}
	return tx
}
