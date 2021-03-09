package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server()

}

const priKey = `CAASqAkwggSkAgEAAoIBAQDRmDUu+idk1SRLKKvjzZ9SAzekFsQuLk/ZNHQ/4uGc6GiFrgjcPvnzuH9MZTHBSZqD5d4c6osWqcjrhzKFOB/EoIekeG9YAl7IZqjNPZEdYlfO04Se60p9PJBU7cinZEw31+aph00f8WECeQuhA3DABk3dDtlSfXb/8iATpUOS001R3dbKXnYMSdzMlzfvrRG9R7GuRyzq1KA8nZoU+tKenNOydhRR4lCKYEcuTkStIfrbU8GOiKyW10GN6UcDkcLqLBJzF3WaPB4g5/1aRSjgrkoUHaOAe/eI6nUnaq0Srvf/Lj5IhzQNpgRuWeBaHAvcxi5T5kXD5jhG0nHrjTbFAgMBAAECggEBAJTMaV17jQIh641MR5QHxAcWb6cX3pkmmduLSMieSyv5N1NIddPfBdbIcd+LrCIcHg1r4R9ISAyD5zeHGQ/JA4y1pvbV5b5mmGHPuCFVhzOEQPB97BZi1tRIbfNNiPmF1DKFHaXXf6Kg3A1EYEQkTCSXlSnOQ+0zU4trmug3mNyfy3jZvG9shp6Ij7aBcLcvkRkzyOkA5WPgiKIeDAzn2gDUOuP4r9qPpWy3Ti9Kax5H+2l/RQ0T0T+WlbDg6yXDRzcMcUb5mWY1XVElVVmiQgl5gu6L6CL2bb/UaGF8ex91Ckse82ckiAShcddaEmYWPN9UqJLp4Z2OlDLfLUin3UECgYEA8huasBuABI3Wf7vThhNsJlMK8Sr4E28NGlAgx0g+Q06ckyJxbVyHJXJ/nktzgEDUXCHzRVNTndSYu/2vRqzNMFNlm7h/YoJghHReMC/rx//ONzKiFr9uFy2sI8BY4T+qKefXdS2Fsk7UJ3t5rzlxLoKu6DImeTZV/zNOa0zfXxUCgYEA3Z8BleXV/0u8bPdnbtchJ2gTWlxPgpIoPNiBY6q9Nw8sSaasSFjYBLasCZKc6NUoYJ+2AheJjvDcEGzFzPSwaxIMSkkUR11vJ6eruDrUkCYtgPOc+jF3Vpo1o88y9cG64nWYdOvt+SLm48touXLjd0AURTF9Cy0omKFQ3Z++5PECgYEAvYkSvo+o1ufbZsA6Rhpqbk5QoKDM+RnVHiZgouJRrAuc1CsAtWbcflp2wgu7bkpSdZY2hq1HJqZKs9FUKHYbZJvFTfVP9GSw/sDDA+JgKYB/hgLjlf9jRk4BFzP74MsgghH4QMnUgtTnjclCaAUGMC0qlKi+KeJ5zIH0AFh7/kkCgYACyHdloYBBd3sDR0wWOT9iVk0/6j7ZXeqBcRqW3NMJePhOaHhrZCo6TOz2JdAwoFSkefz4I8GHeQDad/M38q6weYaL/ETz7Hlz3wgqBRscQE57+xMylSJxhPg9eWGjcm5dX6qtdTUE1updW/WRtp6ipbxbbhaq6ENFP2lbjyD/kQKBgHHF+lT6RJVtrGcfmwJsD5yG2c79lKGaz4UMJShKyD9wb8y1lFEHRbi9HJcJidEFh13GmjOU5QJRXGq3B4grBM3fCmYzbW2dlDcsUQeEcU3GYF7jqSEdfebPRVo+xKEaUFnj+bfbBs0nwMyuEeMaKniMQb0MxGiyiWzLSkcD60Di`

const CUSTOM_ID0 protocol.ID = "from_nemo0"

const CUSTOM_ID1 protocol.ID = "from_free"

const WAIT protocol.ID = "/wait/1.0.0"

func server() {
	addr0, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/63729")
	if err != nil {
		panic(err)
	}
	addr1, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/63730")
	if err != nil {
		panic(err)
	}

	p, _ := base64.RawStdEncoding.DecodeString(priKey)
	c, err := crypto.UnmarshalPrivateKey(p)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	host, err := libp2p.New(ctx,
		libp2p.ListenAddrs(addr0, addr1),
		libp2p.Ping(false),
		libp2p.Identity(c),
	)
	if err != nil {
		panic(err)
	}

	res := host.Peerstore().PrivKey(host.ID())

	a, _ := res.Bytes()

	fmt.Println(base64.RawStdEncoding.EncodeToString(a))

	b, _ := crypto.UnmarshalPrivateKey(a)

	fmt.Println(b.Bytes())

	host.SetStreamHandler(CUSTOM_ID0, Deal0)
	host.SetStreamHandler(CUSTOM_ID1, Deal1)
	host.SetStreamHandler(WAIT, wait)

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

	for _, addr := range addrs {
		fmt.Println("addr:", addr)
	}

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
	for {
		n, err := s.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(buf[0:n]))
	}
}

func Deal1(s network.Stream) {
	fmt.Println("deal1")
	fmt.Println("id:", s.ID())
	fmt.Println("protocol:", s.Protocol())
}

type waitStru struct {
	T    int64
	Info string
}

func wait(s network.Stream) {
	fmt.Println("deal2")
	fmt.Println("id:", s.ID())
	fmt.Println("protocol:", s.Protocol())
	buf := make([]byte, 256)
	n, err := s.Read(buf)
	if err != nil {
		panic(err)
	}
	var rec waitStru
	json.Unmarshal(buf[:n], &rec)
	fmt.Println(rec.Info)
	time.Sleep(time.Duration(rec.T))
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
