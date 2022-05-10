package library

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	noise "github.com/libp2p/go-libp2p-noise"
	"github.com/multiformats/go-multiaddr"
)

const ip = "127.0.0.1"

type NodeConfig struct {
	Port          int
	PrivateKey    libp2pcrypto.PrivKey
	CloseInbound  bool
	CloseOutbound bool
}

type node struct {
	libp2pHost host.Host
	port       int
}

func NewNode(config NodeConfig) (*node, error) {
	listenAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ip, config.Port))
	if err != nil {
		return nil, fmt.Errorf("libp2p multiaddress error - can not start host on port %d", config.Port)
	}

	libp2pHost, err := libp2p.New(
		libp2p.Security(noise.ID, noise.New),
		libp2p.ListenAddrs(listenAddr),
		libp2p.NATPortMap(),
		libp2p.Identity(config.PrivateKey),
		// libp2p.ConnectionGater(srv),
	)
	if err != nil {
		return nil, fmt.Errorf("libp2p error - can not start host on port %d", config.Port)
	}

	libp2pHost.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(n network.Network, conn network.Conn) {
			Log("ConnectedF peer = %s, direction = %s, connId = %s",
				conn.RemotePeer(), conn.Stat().Direction, conn.ID())

			if (conn.Stat().Direction == network.DirInbound && config.CloseInbound) ||
				(conn.Stat().Direction == network.DirOutbound && config.CloseOutbound) {
				Log("ConnectedF closing the connection - peer = %s, direction = %s, connId = %s",
					conn.RemotePeer(), conn.Stat().Direction, conn.ID())
				go func() {
					time.Sleep(time.Millisecond * 1000)
					conn.Close()
				}()
			}
		},
		DisconnectedF: func(n network.Network, conn network.Conn) {
			Log("DisconnectedF peer = %s, direction = %s, connId = %s",
				conn.RemotePeer(), conn.Stat().Direction, conn.ID())
		},
	})

	return &node{port: config.Port, libp2pHost: libp2pHost}, nil
}

func (n node) PeerId() peer.ID {
	return n.libp2pHost.ID()
}

func (n *node) Dial(dstMultiAddress string) error {
	addr, err := multiaddr.NewMultiaddr(dstMultiAddress)
	if err != nil {
		return fmt.Errorf("error while dialing peer - NewMultiaddr: %s, err: %v", dstMultiAddress, err)
	}
	peer, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return fmt.Errorf("error while dialing peer - AddrInfoFromP2pAddr: %s, err: %v", dstMultiAddress, err)
	}
	err = n.libp2pHost.Connect(context.Background(), *peer)
	if err != nil {
		return fmt.Errorf("error while dialing peer - Connect: %s, err: %v", dstMultiAddress, err)
	}
	return nil
}

func (n *node) Address() string {
	ma := fmt.Sprintf("/ip4/%s/tcp/%d/p2p/%s", ip, n.port, n.libp2pHost.ID())
	return ma
}

func (n *node) Stop() {
	n.libp2pHost.Close()
}
