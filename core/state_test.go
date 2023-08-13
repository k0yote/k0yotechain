package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {

	state := NewState()
	assert.True(t, len(state.data) == 0)
	assert.Nil(t, state.Put([]byte("key"), []byte("value")))
	value, err := state.Get([]byte("key"))
	assert.Nil(t, err)
	assert.Equal(t, []byte("value"), value)
	assert.Nil(t, state.Delete([]byte("key")))
	_, err = state.Get([]byte("key"))
	assert.NotNil(t, err)
}
