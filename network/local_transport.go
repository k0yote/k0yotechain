package network

import (
	"bytes"
	"fmt"
	"net"
	"sync"
)

type LocalTransport struct {
	addr     net.Addr
	consumCh chan RPC
	lock     sync.RWMutex
	peers    map[net.Addr]*LocalTransport
}

func NewLocalTransport(addr net.Addr) *LocalTransport {
	return &LocalTransport{
		addr:     addr,
		consumCh: make(chan RPC, 1024),
		peers:    make(map[net.Addr]*LocalTransport),
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

func (t *LocalTransport) SendMessage(to net.Addr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if t.addr == to {
		return nil
	}

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send message to unkown peer %s", t.addr, to)
	}

	peer.consumCh <- RPC{
		From:    t.addr,
		Payload: bytes.NewReader(payload),
	}

	return nil
}

func (t *LocalTransport) Broadcast(payload []byte) error {
	for _, peer := range t.peers {
		if err := t.SendMessage(peer.Addr(), payload); err != nil {
			return err
		}
	}
	return nil
}

func (t *LocalTransport) Addr() net.Addr {
	return t.addr
}
