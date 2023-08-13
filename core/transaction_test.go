package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/k0yote/privatechain/crypto"
	"github.com/k0yote/privatechain/types"
	"github.com/stretchr/testify/assert"
)

func TestVerifyTransactionWithTamper(t *testing.T) {
	tx := NewTransaction(nil)
	from := crypto.GeneratePrivateKey()
	to := crypto.GeneratePrivateKey()
	hacker := crypto.GeneratePrivateKey()

	tx.From = from.PublicKey()
	tx.To = to.PublicKey()
	tx.Value = 1_000
	assert.Nil(t, tx.Sign(from))
	tx.hash = types.Hash{}
	tx.To = hacker.PublicKey()

	assert.NotNil(t, tx.Verify())
}

func TestNativeTransferTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	toPrivKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		To:    toPrivKey.PublicKey(),
		Value: 10,
	}

	hash := tx.Hash(TxHasher{})
	fmt.Println(hash)

	assert.Nil(t, tx.Sign(fromPrivKey))
}

func TestNFTTransaction(t *testing.T) {

	collectionTx := CollectionTx{
		Fee:      10,
		MetaData: []byte("the beginning of a new collection"),
	}

	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		TxInner: collectionTx,
	}

	tx.Sign(privKey)
	buf := new(bytes.Buffer)
	assert.Nil(t, gob.NewEncoder(buf).Encode(tx))
	tx.hash = types.Hash{}

	txDecoded := &Transaction{}
	assert.Nil(t, gob.NewDecoder(buf).Decode(txDecoded))
	assert.Equal(t, tx, txDecoded)
}

func TestSignTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.Nil(t, tx.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	tx.From = otherPrivKey.PublicKey()

	assert.NotNil(t, tx.Verify())
}

func TestTxEncodeDecode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t, tx.Encode(NewGobTxEncoder(buf)))
	tx.hash = types.Hash{}

	txDecoded := new(Transaction)
	assert.Nil(t, txDecoded.Decode(NewGobTxDecoder(buf)))
	assert.Equal(t, tx, txDecoded)
}

func randomTxWithSignature(t *testing.T) *Transaction {
	privKey := crypto.GeneratePrivateKey()
	tx := Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	return &tx
}
