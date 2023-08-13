package core

import (
	"testing"

	"github.com/k0yote/privatechain/crypto"
	"github.com/stretchr/testify/assert"
)

func TestAccountState(t *testing.T) {
	state := NewAccountState()

	address := crypto.GeneratePrivateKey().PublicKey().Address()

	account := state.CreateAccount(address)

	assert.Equal(t, address, account.Address)
	assert.Equal(t, uint64(0), account.Balance)
	got, err := state.GetAccount(address)
	assert.Nil(t, err)
	assert.Equal(t, account, got)

	addressOther := crypto.GeneratePrivateKey().PublicKey().Address()

	accountOther := state.CreateAccountWithBalance(addressOther, 1_000_000)
	assert.Equal(t, addressOther, accountOther.Address)
	assert.Equal(t, uint64(1_000_000), accountOther.Balance)
	gotOther, err := state.GetAccount(addressOther)
	assert.Nil(t, err)
	assert.Equal(t, accountOther, gotOther)

	assert.Nil(t, state.Transfer(addressOther, address, 1_000))

	addressBalance, err := state.GetBalance(address)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1_000), addressBalance)

	addressOtherBalance, err := state.GetBalance(addressOther)
	assert.Nil(t, err)
	assert.Equal(t, uint64(999_000), addressOtherBalance)

	aliceAddress := crypto.GeneratePrivateKey().PublicKey().Address()
	_, err = state.GetAccount(aliceAddress)

	assert.EqualError(t, ErrAccountNotFound, err.Error())

	_, err = state.GetBalance(aliceAddress)
	assert.EqualError(t, ErrAccountNotFound, err.Error())

	assert.EqualError(t, ErrAccountNotFound, state.Transfer(aliceAddress, address, 2000).Error())

	bobAddress := crypto.GeneratePrivateKey().PublicKey().Address()
	bobAccount := state.CreateAccountWithBalance(bobAddress, 2000)

	assert.Equal(t, "2000", bobAccount.String())

	etcAddress := crypto.GeneratePrivateKey().PublicKey().Address()
	assert.EqualError(t, ErrInsufficientBalance, state.Transfer(bobAddress, etcAddress, 3000).Error())

	assert.Nil(t, state.Transfer(bobAddress, etcAddress, 2000))

	etcBalance, err := state.GetBalance(etcAddress)
	assert.Nil(t, err)
	assert.Equal(t, uint64(2000), etcBalance)
	bobBalance, err := state.GetBalance(bobAddress)
	assert.Nil(t, err)
	assert.Equal(t, uint64(0), bobBalance)
}
