package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server()

}

const CUSTOM_ID0 protocol.ID = "from_nemo0"

const CUSTOM_ID1 protocol.ID = "from_free"

func server() {
	addr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/63729")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	host, err := libp2p.New(ctx,
		libp2p.ListenAddrs(addr),
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}
	host.SetStreamHandler(CUSTOM_ID0, Deal0)
	host.SetStreamHandler(CUSTOM_ID1, Deal1)

	peerInfo := peerstore.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("host_id:", host.ID())
	fmt.Println("host_addrs:", host.Addrs())
	fmt.Println("addr:", addrs[0])

	exit_ch := make(chan os.Signal, 1)
	signal.Notify(exit_ch, syscall.SIGINT, syscall.SIGTERM)
	<-exit_ch
	fmt.Println("Received signal, shutting down...")
	if err := host.Close(); err != nil {
		panic(err)
	}
}

func Deal0(s network.Stream) {
	fmt.Println("deal0")
	fmt.Println("id:", s.ID())
	fmt.Println("protocol:", s.Protocol())

	buf := make([]byte, 256)
	n, err := s.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf[0:n]))
}

func Deal1(s network.Stream) {
	fmt.Println("deal1")
	fmt.Println("id:", s.ID())
	fmt.Println("protocol:", s.Protocol())
}

func Deal2(s network.Stream) {
	fmt.Println("deal2")
	fmt.Println("id:", s.ID())
	fmt.Println("protocol:", s.Protocol())
}

func serverPing() {
	// create a background context (i.e. one that never cancels)
	ctx := context.Background()

	// start a libp2p node that listens on a random local TCP port,
	// but without running the built-in ping protocol
	node, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}

	// configure our own ping protocol
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, Hello(pingService.PingHandler))

	// print the node's PeerInfo in multiaddr format
	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("libp2p node address:", addrs[0])
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
	if err := node.Close(); err != nil {
		panic(err)
	}
}

func Hello(do func(s network.Stream)) func(s network.Stream) {
	return func(s network.Stream) {
		fmt.Println("id:", s.ID(), "protocol:", s.Protocol(), "stat:", s.Stat())
		do(s)
	}
}
