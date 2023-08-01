package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeypairSignVerifySuccess(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()

	msg := []byte("Hello World")
	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)
	assert.True(t, sig.Verify(pubKey, msg))
}

func TestKeypairSignVerifyFailure(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()

	msg := []byte("Hello World")
	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)
	assert.False(t, sig.Verify(pubKey, []byte("aaa")))
}
