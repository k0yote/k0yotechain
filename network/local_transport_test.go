package network

import (
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dummy, _ = net.ResolveTCPAddr("tcp", "A")
)

func TestConnect(t *testing.T) {

	tra := NewLocalTransport(dummy)
	trb := NewLocalTransport(dummy)

	tra.Connect(trb)
	trb.Connect(tra)
	assert.Equal(t, tra.peers[trb.addr], trb)
	assert.Equal(t, trb.peers[tra.addr], tra)
}

func TestSendMessge(t *testing.T) {
	tra := NewLocalTransport(dummy)
	trb := NewLocalTransport(dummy)

	tra.Connect(trb)
	trb.Connect(tra)

	msg := []byte("hello world")
	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))

	rpc := <-trb.Consume()
	b, err := ioutil.ReadAll(rpc.Payload)
	assert.Nil(t, err)

	assert.Equal(t, b, msg)
	assert.Equal(t, rpc.From, tra.Addr())
}

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport(dummy)
	trb := NewLocalTransport(dummy)
	trc := NewLocalTransport(dummy)

	tra.Connect(trc)
	tra.Connect(trb)

	msg := []byte("foo")
	assert.Nil(t, tra.Broadcast(msg))

	rpcb := <-trb.Consume()
	b, err := ioutil.ReadAll(rpcb.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)

	rpcc := <-trc.Consume()
	b, err = ioutil.ReadAll(rpcc.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)
}
