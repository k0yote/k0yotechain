package network

type GetBlocksMessage struct {
	From uint32
	To   uint32
}

type GetStatusMessage struct{}

type StatusMessage struct {
	ID            string
	Version       uint32
	CurrentHeight uint32
}
