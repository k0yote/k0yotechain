package network

import (
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr     NetAddr
	consumCh chan RPC
	lock     sync.RWMutex
	peers    map[NetAddr]*LocalTransport
}

func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:     addr,
		consumCh: make(chan RPC, 1024),
		peers:    make(map[NetAddr]*LocalTransport),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumCh
}

func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[tr.Addr()] = tr.(*LocalTransport)

	return nil
}

func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send message to %s", t.addr, to)
	}

	peer.consumCh <- RPC{
		From:    t.addr,
		Payload: payload,
	}

	return nil
}

func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}
