package udp_client

import (
	"encoding/base64"
	"net"
	"time"

	"p2p-network-test/base_client"
	"p2p-network-test/udp_server"
)

type Client struct {
	*base_client.Client
	sAddr *net.UDPAddr
}

func (c *Client) Start() error {
	s := c.GetServer()

	sConn, err := s.CreateConn(c.Addr)
	if err != nil {
		return err
	}

	//  pubKey, err := c.GetSelf().GetPublicKey()
	// 	if err != nil {
	// 		return err
	// 	}

	go s.Listen()

	sConn.Send(&shared.Message{
		Type:    "greeting",
		Content: "hello",
	})

	return nil
}

func (c *Client) Connect() {
	l := c.GetLog()
	self := c.GetSelf()
	peer := c.GetPeer()
	pConn := c.GetPeerConn()

	for i := 0; i < 5; i++ {
		if c.WasKeyReceived() {
			l.Printf("connected to peer %s", peer.Username)
			c.ConnectedCallback(c)
			return
		}

		l.Printf("punching through to peer %s at %s", peer.Username, pConn.GetAddr())
		pConn.Send(&shared.Message{
			Type:   "connect",
			PeerID: self.ID,
		})
		time.Sleep(3 * time.Second)
	}

	l.Printf("could not connect to peer %s at %s", peer.Username, pConn.GetAddr())
}

func New(username string, addr *net.UDPAddr) {
	//创建udp服务器
	s, err := udp_server.New(addr)
	if err != nil {
		return nil, err
	}

	bc, err := base_client.New(username, s)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Client: bc,
		sAddr:  sAddr,
	}

	s.OnMessage(shared.CreateMessageCallback(c))

	return c, nil
}
