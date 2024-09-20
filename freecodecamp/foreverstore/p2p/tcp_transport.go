package p2p

import "net"

type TCPTransport struct {
	ListenAddress string
	Listener      net.Listener
	peers         map[string]Peer
}
