package network

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/k0yote/privatechain/core"
	"github.com/k0yote/privatechain/crypto"
	"github.com/k0yote/privatechain/types"
	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	SeedNodes     []string
	ListenAddr    string
	TCPTransport  *TCPTransport
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
	Blockchain    *core.Blockchain
}

type Server struct {
	TCPTransport *TCPTransport
	peerCh       chan *TCPPeer
	peerMap      map[net.Addr]*TCPPeer
	ServerOpts
	mempool     *TxPool
	chain       *core.Blockchain
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "addr", opts.ID)
	}

	chain, err := core.NewBlockchain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}

	peerCh := make(chan *TCPPeer)
	tr := NewTCPTransport(opts.ListenAddr, peerCh)

	s := &Server{
		TCPTransport: tr,
		peerCh:       peerCh,
		peerMap:      make(map[net.Addr]*TCPPeer),
		ServerOpts:   opts,
		mempool:      NewTxPool(1000),
		chain:        chain,
		isValidator:  opts.PrivateKey != nil,
		rpcCh:        make(chan RPC),
		quitCh:       make(chan struct{}, 1),
	}

	s.TCPTransport.peerCh = peerCh

	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}

func (s *Server) bootstrapNetwork() {
	for _, addr := range s.SeedNodes {
		go func(addr string) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Printf("could not connect to %+v [err: %s]\n", conn, err)
				return
			}

			s.peerCh <- &TCPPeer{
				conn: conn,
			}
		}(addr)
	}
}

func (s *Server) Start() {
	s.TCPTransport.Start()

	time.Sleep(1 * time.Second)

	s.bootstrapNetwork()

	s.Logger.Log("msg", "accepting TCP connection on", "addr", s.ListenAddr, "id", s.ID)

free:
	for {
		select {
		case peer := <-s.peerCh:
			s.peerMap[peer.conn.RemoteAddr()] = peer

			go peer.readLoop(s.rpcCh)
			fmt.Printf("New peer: %+v\n", peer)
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)
				continue
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				if err != core.ErrBlockKnown {
					logrus.Error(err)
				}
			}

		case <-s.quitCh:
			break free
		}
	}

	s.Logger.Log("msg", "Server is shutting down")
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)

	s.Logger.Log("msg", "Starting validator loop", "BlockTime", s.BlockTime)

	for {
		<-ticker.C
		s.createNewBlock()
	}
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	case *core.Block:
		return s.processBlock(t)
	case *GetStatusMessage:
		// return s.processGetStatusMessage(msg.From, t)
	case *StatusMessage:
		// return s.processStatusMessage(msg.From, t)
	case *GetBlocksMessage:
		return s.processGetBlocksMessage(msg.From, t)
	}

	return nil
}

func (s *Server) processGetBlocksMessage(from net.Addr, data *GetBlocksMessage) error {
	fmt.Printf("%+v\n", data)
	return nil
}

// func (s *Server) sendGetStatusMessage(tr Transport) error {
// 	var (
// 		getStatausMsg = new(GetStatusMessage)
// 		buf           = new(bytes.Buffer)
// 	)

// 	if err := gob.NewEncoder(buf).Encode(getStatausMsg); err != nil {
// 		return nil
// 	}

// 	msg := NewMessage(MessageTypeGetStatus, buf.Bytes())

// 	if err := s.Transport.SendMessage(tr.Addr(), msg.Bytes()); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (s *Server) broadcast(payload []byte) error {
	fmt.Println("broadcasting: ")
	for netAddr, peer := range s.peerMap {
		if err := peer.Send(payload); err != nil {
			fmt.Printf("peer send error => addr %s [err: %s]", netAddr, err)
		}
	}

	return nil
}

// func (s *Server) processStatusMessage(from NetAddr, data *StatusMessage) error {
// 	if data.CurrentHeight <= s.chain.Height() {
// 		s.Logger.Log("msg", "cannot sync blockHeight to low", "ourHeight", s.chain.Height(), "theirHeight", data.CurrentHeight, "addr", from)
// 		return nil
// 	}

// 	getBlocksMessage := &GetBlocksMessage{
// 		From: s.chain.Height(),
// 		To:   0,
// 	}

// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(getBlocksMessage); err != nil {
// 		return err
// 	}

// 	msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())

// 	return s.Transport.SendMessage(from, msg.Bytes())
// }

// func (s *Server) processGetStatusMessage(from NetAddr, gmsg *GetStatusMessage) error {
// 	fmt.Printf("received status from %s => %+v\n", from, gmsg)

// 	statusMessage := &StatusMessage{
// 		CurrentHeight: s.chain.Height(),
// 		ID:            s.ID,
// 	}

// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(statusMessage); err != nil {
// 		return err
// 	}

// 	msg := NewMessage(MessageTypeStatus, buf.Bytes())

// 	return s.Transport.SendMessage(from, msg.Bytes())
// }

func (s *Server) processBlock(b *core.Block) error {
	if err := s.chain.AddBlock(b); err != nil {
		return err
	}

	go s.broadcastBlock(b)
	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.mempool.Contains(hash) {

		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	// s.Logger.Log(
	// 	"msg", "adding new tx to mempool",
	// 	"hash", hash,
	// 	"mempoolPending", s.mempool.PendingCount(),
	// )

	go s.broadcastTx(tx)

	s.mempool.Add(tx)

	return nil
}

func (s *Server) broadcastBlock(b *core.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(core.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeBlock, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	txx := s.mempool.Pending()
	block, err := core.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return err
	}

	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	s.mempool.ClearPending()

	go s.broadcastBlock(block)

	return nil
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Height:    0,
		Timestamp: 000000,
	}

	b, _ := core.NewBlock(header, nil)
	return b
}
