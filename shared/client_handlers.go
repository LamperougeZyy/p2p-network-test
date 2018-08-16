package shared

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"

	"github.com/mitchellh/mapstructure"
)

func greetingHandler(c Client, serverConn Conn, m *Message) (*Message, error) {
	l := c.GetLog()
	self := c.GetSelf()

	if m.Error != "" {
		l.Fatal(m.Error)
		return nil, errors.New(m.Error)
	}

	s, ok := m.Content.(string)
	if !ok {
		return nil, errors.New("expected to received public key with greeting")
	}

	return &Message{
		Type:   "register",
		PeerID: self.ID,
		Content: Registration{
			Username:  self.Username,
			PublicKey: "111",
		},
	}, nil
}

func registerHandler(c Client, serverConn Conn, m *Message) (*Message, error) {
	if m.Error != "" {
		return nil, errors.New(m.Error)
	}

	c.RegisteredCallback(c)
	return nil, nil
}

func establishHandler(c Client, serverConn Conn, m *Message) (*Message, error) {
	l := c.GetLog()
	l.Print("establish request from server")
	if m.Error != "" {
		return nil, errors.New(m.Error)
	}

	var p Peer
	err := mapstructure.Decode(m.Content, &p)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	c.SetPeer(&Peer{
		ID:       p.ID,
		Username: p.Username,
	})

	var addr net.Addr
	switch serverConn.Protocol() {
	case "UDP":
		addr, err = net.ResolveUDPAddr("udp", p.Endpoint.String())
	case "TCP":
		addr, err = net.ResolveTCPAddr("tcp", p.Endpoint.String())
	default:
		addr, err = nil, fmt.Errorf("unknown Conn protocol %s", serverConn.Protocol())
	}
	if err != nil {
		return nil, err
	}

	if c.GetPeerConn() != nil && c.GetPeerConn().GetAddr().String() != addr.String() {
		l.Print("ignoring establish request because the client is already connected to a peer")
		return nil, nil
	}

	go func() {
		pConn, err := c.GetServer().CreateConn(addr)
		if err != nil {
			return
		}

		c.SetPeerConn(pConn)

		go c.Connect()

		c.ConnectingCallback(c)
	}()
	return nil, nil
}

func connectHandler(c Client, peerConn Conn, m *Message) (*Message, error) {
	self := c.GetSelf()
	l := c.GetLog()

	pConn := c.GetPeerConn()
	if pConn == nil {
		return nil, nil
	}

	if pConn != peerConn {
		//如果收到了一样的地址，那么监听者就需要建立好对应的连接
		if pConn.GetAddr().String() == peerConn.GetAddr().String() {
			pConn = peerConn
			c.SetPeerConn(pConn)
		}
	}

	l.Printf("connection mirror request from peer %s at %s, sending mirror", self.Username, pConn.GetAddr())
	pubKey, err := self.GetPublicKey()
	if err != nil {
		return nil, err
	}

	defer c.SetKeySent(true)

	return &Message{
		Type:    "key",
		PeerID:  self.ID,
		Content: "2222",
	}, nil
}
