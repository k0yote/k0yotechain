package network

import "github.com/k0yote/privatechain/core"

type GetBlocksMessage struct {
	From uint32
	To   uint32
}

type BlocksMessage struct {
	Blocks []*core.Block
}

type GetStatusMessage struct{}

type StatusMessage struct {
	ID            string
	Version       uint32
	CurrentHeight uint32
}
