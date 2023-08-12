package core

import (
	"testing"

	"github.com/k0yote/privatechain/crypto"
	"github.com/stretchr/testify/assert"
)

func TestAccountStateTransferNoBalance(t *testing.T) {
	state := NewAccountState()

	from := crypto.GeneratePrivateKey().PublicKey().Address()
	to := crypto.GeneratePrivateKey().PublicKey().Address()
	amount := uint64(10)

	assert.NotNil(t, state.Transfer(from, to, amount))
}

func TestAccountStateTransferSuccess(t *testing.T) {
	state := NewAccountState()
	from := crypto.GeneratePrivateKey().PublicKey().Address()

	to := crypto.GeneratePrivateKey().PublicKey().Address()
	amount := uint64(10)

	assert.Nil(t, state.Transfer(from, to, amount))
}
