package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)
	assert.Equal(t, tra.Addr(), NetAddr("A"))
	assert.Equal(t, trb.Addr(), NetAddr("B"))
}

func TestSendMessge(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	msg := []byte("hello world")
	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))

	rpc := <-trb.Consume()
	assert.Equal(t, rpc.Payload, msg)
	assert.Equal(t, rpc.From, tra.Addr())
}
