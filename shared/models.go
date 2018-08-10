package shared

import (
	"encoding/base64"
	"errors"
	"net"
	"strconv"
)

type Conn interface {
	Send(*Message) error
	Protocol() string
	GetAddr() net.Addr
	GetSecret() ([32]byte, error)
	SetSecret([32]byte)
}

type UDPPayload struct {
	Bytes []byte
	Addr  *net.UDPAddr
}

type UDPConn struct {
	send   chan *UDPPayload
	addr   *net.UDPAddr
	secret string
}

func convertSecret(secretText string) ([32]byte, error) {
	var secret [32]byte
	if secretText == "" {
		return secret, errors.New("secret has not been set")
	}

	//对秘文进行解密，结果存储到byte切片里
	bs, err := base64.StdEncoding.DecodeString(secretText)
	if err != nil {
		return secret, errors.New("could not decode secret")
	}

	copy(secret[:], bs)
	return secret, nil
}

func (c *UDPConn) GetSecret() ([32]byte, error) {
	return convertSecret(c.secret)
}

func (c *UDPConn) SetSecret([32]byte) {
}

func (c *UDPConn) Send(m *Message) error {
	//首先要对massage进行解析，得到json格式的切片
	b, err := MessageOut(c, m)
	if err != nil {
		return err
	}

	c.send <- &UDPPayload{Bytes: b, Addr: c.addr}
	return err
}

func (c *UDPConn) GetAddr() net.Addr {
	return c.addr
}

func (c *UDPConn) Protocol() string {
	return "UDP"
}

func NewUDPConn(send chan *UDPPayload, addr *net.UDPAddr) *UDPConn {
	return &UDPConn{
		send: send,
		addr: addr,
	}
}

type Conns map[string]Conn

type Registration struct {
	Username string `json:"username"`
	//PublicKey string `json:"publickey"`
}

type Message struct {
	Type    string      `json:"type"`
	PeerID  string      `json:"peerID, omitempty"`
	Error   string      `json:"error, omitempty"`
	Content interface{} `json:"data, omitempty"`
	Encrypt bool        `json:"-"`
	addr    *net.UDPAddr
}

type Endpoint struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

func (e Endpoint) String() string {
	return e.IP + ":" + strconv.Itoa(e.Port)
}

type Peer struct {
	ID         string       `json:"id, omitempty"`
	Username   string       `json:"username, omitempty"`
	Endpoint   Endpoint     `json:"endpoint,omitempty"`
	PublicKey  string       `json:"publickey, omitempty"`
	PrivateKey [32]byte     `json:"-"`
	Addr       *net.UDPAddr `json:"-"`
}

type Peers map[string]*Peer
