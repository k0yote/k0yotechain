package main

import (
	"bytes"
	"log"
	"net"

	"github.com/k0yote/privatechain/core"
	"github.com/k0yote/privatechain/crypto"
	"github.com/k0yote/privatechain/network"
)

func main() {
	privKey := crypto.GeneratePrivateKey()
	localNode := makeServer("LOCAL_NODE", &privKey, ":3000", []string{":4000"})

	go localNode.Start()

	remoteNodeA := makeServer("REMOTE_NODE_A", nil, ":4000", []string{":4001"})
	go remoteNodeA.Start()

	remoteNodeB := makeServer("REMOTE_NODE_B", nil, ":4001", nil)
	go remoteNodeB.Start()
	// time.Sleep(1 * time.Second)

	// tcpTester()

	select {}
}

func makeServer(id string, pk *crypto.PrivateKey, addr string, seedNodes []string) *network.Server {
	opts := network.ServerOpts{
		SeedNodes:  seedNodes,
		ListenAddr: addr,
		PrivateKey: pk,
		ID:         id,
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func tcpTester() {
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		panic(err)
	}

	privkey := crypto.GeneratePrivateKey()

	tx := core.NewTransaction(contract())
	tx.Sign(privkey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		log.Fatal(err)
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	_, err = conn.Write(msg.Bytes())
	if err != nil {
		panic(err)
	}
}

// var (
// 	transports = []network.Transport{
// 		network.NewLocalTransport("LOCAL"),
// 		network.NewLocalTransport("REMOTE_A"),
// 		network.NewLocalTransport("REMOTE_B"),
// 		network.NewLocalTransport("REMOTE_C"),
// 		network.NewLocalTransport("LATE_REMOTE"),
// 	}
// )

// func main() {
// 	initRemoteServers(transports)

// 	localNode := transports[0]
// 	remoteNodeA := transports[1]
// 	remoteNodeC := transports[3]

// 	go func() {
// 		for {
// 			if err := sendTransaction(remoteNodeA, localNode.Addr()); err != nil {
// 				logrus.Error(err)
// 			}
// 			time.Sleep(2 * time.Second)
// 		}
// 	}()

// 	go func() {
// 		time.Sleep(7 * time.Second)
// 		trLate := transports[len(transports)-1]
// 		remoteNodeC.Connect(trLate)
// 		lateServer := makeServer(string(trLate.Addr()), trLate, nil)
// 		go lateServer.Start()
// 	}()

// 	privKey := crypto.GeneratePrivateKey()

// 	localServer := makeServer("LOCAL", transports[0], &privKey)
// 	localServer.Start()
// }

// func initRemoteServers(trs []network.Transport) {
// 	for i := 0; i < len(trs); i++ {
// 		id := fmt.Sprintf("REMOTE_%d", i+1)
// 		s := makeServer(id, trs[i], nil)
// 		go s.Start()
// 	}
// }

// func sendGetStatusMessage(tr network.Transport, to network.NetAddr) error {
// 	var (
// 		getStatausMsg = new(network.GetStatusMessage)
// 		buf           = new(bytes.Buffer)
// 	)

// 	if err := gob.NewEncoder(buf).Encode(getStatausMsg); err != nil {
// 		return nil
// 	}

// 	msg := network.NewMessage(network.MessageTypeGetStatus, buf.Bytes())

// 	return tr.SendMessage(to, msg.Bytes())
// }

// func sendTransaction(tr network.Transport, to network.NetAddr) error {
// 	privkey := crypto.GeneratePrivateKey()

// 	tx := core.NewTransaction(contract())
// 	tx.Sign(privkey)
// 	buf := &bytes.Buffer{}
// 	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
// 		return err
// 	}

// 	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

// 	return tr.SendMessage(to, msg.Bytes())
// }

func contract() []byte {
	data := []byte{0x02, 0x0a, 0x03, 0x0a, 0x0b, 0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x0f}
	pushFoo := []byte{0x4f, 0x0c, 0x4f, 0x0c, 0x46, 0x0c, 0x03, 0x0a, 0x0d, 0x10}
	data = append(data, pushFoo...)

	return data
}
