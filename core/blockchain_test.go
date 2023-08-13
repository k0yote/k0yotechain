package core

import (
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/k0yote/privatechain/crypto"
	"github.com/k0yote/privatechain/types"
	"github.com/stretchr/testify/assert"
)

func TestSendNativeTransferTampered(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	signer := crypto.GeneratePrivateKey()
	block := randomBlock(t, uint32(1), getPrevBlockHash(t, bc, uint32(1)))
	assert.Nil(t, block.Sign(signer))

	bob := crypto.GeneratePrivateKey()
	alice := crypto.GeneratePrivateKey()

	bc.accountState.CreateAccountWithBalance(alice.PublicKey().Address(), 1_000_000)

	tx := NewTransaction(nil)
	tx.From = alice.PublicKey()
	tx.To = bob.PublicKey()
	tx.Value = 1_000

	assert.Nil(t, tx.Sign(alice))
	tx.hash = types.Hash{}

	hacker := crypto.GeneratePrivateKey()
	tx.To = hacker.PublicKey()

	block.AddTransaction(tx)
	assert.NotNil(t, bc.AddBlock(block))

	_, err := bc.accountState.GetAccount(bob.PublicKey().Address())
	assert.EqualError(t, ErrAccountNotFound, err.Error())
}

func TestSendNativeTransferSuccess(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	signer := crypto.GeneratePrivateKey()
	block := randomBlock(t, uint32(1), getPrevBlockHash(t, bc, uint32(1)))
	assert.Nil(t, block.Sign(signer))

	bob := crypto.GeneratePrivateKey()
	alice := crypto.GeneratePrivateKey()

	bc.accountState.CreateAccountWithBalance(alice.PublicKey().Address(), 1_000_000)

	tx := NewTransaction(nil)
	tx.From = alice.PublicKey()
	tx.To = bob.PublicKey()
	tx.Value = 1_000

	assert.Nil(t, tx.Sign(alice))
	block.AddTransaction(tx)
	assert.Nil(t, bc.AddBlock(block))

	balance, err := bc.accountState.GetBalance(bob.PublicKey().Address())
	assert.Nil(t, err)
	assert.Equal(t, uint64(1_000), balance)
}

func TestSendNativeTransferFailNotFound(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	signer := crypto.GeneratePrivateKey()
	block := randomBlock(t, uint32(1), getPrevBlockHash(t, bc, uint32(1)))
	assert.Nil(t, block.Sign(signer))

	bob := crypto.GeneratePrivateKey()
	alice := crypto.GeneratePrivateKey()

	// bc.accountState.CreateAccountWithBalance(alice.PublicKey().Address(), 1_000_000)

	tx := NewTransaction(nil)
	tx.From = alice.PublicKey()
	tx.To = bob.PublicKey()
	tx.Value = 2_000

	assert.Nil(t, tx.Sign(alice))
	block.AddTransaction(tx)
	assert.Nil(t, bc.AddBlock(block))

	hash := tx.Hash(TxHasher{})
	_, err := bc.GetTxByHash(hash)
	assert.NotNil(t, err)
}

func TestSendNativeTransferFailInsufficient(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	signer := crypto.GeneratePrivateKey()
	block := randomBlock(t, uint32(1), getPrevBlockHash(t, bc, uint32(1)))
	assert.Nil(t, block.Sign(signer))

	bob := crypto.GeneratePrivateKey()
	alice := crypto.GeneratePrivateKey()

	bc.accountState.CreateAccountWithBalance(bob.PublicKey().Address(), 1_000)
	tx := NewTransaction(nil)
	tx.From = bob.PublicKey()
	tx.To = alice.PublicKey()
	tx.Value = 2_000

	assert.Nil(t, tx.Sign(bob))
	tx.hash = types.Hash{}
	block.AddTransaction(tx)
	assert.Nil(t, bc.AddBlock(block))

	hash := tx.Hash(TxHasher{})
	_, err := bc.GetTxByHash(hash)
	assert.NotNil(t, err)
}

func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 89, types.Hash{})))
}

func TestNewBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(1))
	assert.False(t, bc.HasBlock(100))
}

func TestGetBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlocks := 100
	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		fetchBlock, err := bc.GetBlock(uint32(i + 1))
		assert.Nil(t, err)
		assert.Equal(t, fetchBlock, block)
	}
}

func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(uint32(i + 1))
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}

func TestAddBlockToHeigh(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.Nil(t, bc.AddBlock(randomBlock(t, 1, getPrevBlockHash(t, bc, uint32(1)))))
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 3, types.Hash{})))
}

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(log.NewLogfmtLogger(os.Stderr), randomBlock(t, 0, types.Hash{}))
	assert.Nil(t, err)
	return bc
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)

	return BlockHasher{}.Hash(prevHeader)
}
