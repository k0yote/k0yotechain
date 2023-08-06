package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := NewStack(128)
	s.Push(1)
	s.Push("a")

	assert.Equal(t, 1, s.Pop())

	assert.Equal(t, "a", s.Pop())
}

func TestStackBytes(t *testing.T) {
	s := NewStack(128)
	s.Push(2)
	s.Push(0x61)
	s.Push(0x62)
}

func TestVM(t *testing.T) {
	// data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}
	pushFoo := []byte{0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x10}
	data = append(data, pushFoo...)

	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())
	// fmt.Printf("%+v\n", vm.stack.data)
	value := vm.stack.Pop().([]byte)
	valueDeSerialize := deSerializeInt64(value)

	assert.Equal(t, valueDeSerialize, int64(5))

	// valueBytes, err := contractState.Get([]byte("FOO"))
	// assert.Nil(t, err)

	// value := deSerializeInt64(valueBytes)
	// assert.Nil(t, err)
	// assert.Equal(t, value, int64(5))
}

func TestVMMul(t *testing.T) {
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x11}

	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())

	result := vm.stack.Pop()
	assert.Equal(t, 6, result)
}

func TestVMDiv(t *testing.T) {
	data := []byte{0x09, 0x0a, 0x03, 0x0a, 0x12}

	contractState := NewState()
	vm := NewVM(data, contractState)
	assert.Nil(t, vm.Run())

	result := vm.stack.Pop()
	assert.Equal(t, 3, result)
}
