package shared

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Conn interface {
	Send(*Message) error
	Protocol() string
	GetAddr() net.Addr
	GetSecret() ([32]byte, error)
	SetSecret([32]byte)
}

type Server interface {
	Stop()
	Listen()
	CreateConn(net.Addr) (Conn, error)
	OnMessage(f func(Conns, Conn, *Message))
}

type Client interface {
	WasKeySent() bool
	SetKeySent(bool)
	WasKeyReceived() bool
	SetKeyReceived(bool)
	GetServer() Server
	GetLog() *log.Logger
	GetSelf() *Peer
	GetPeer() *Peer
	SetPeer(*Peer)
	GetPeerConn() Conn
	SetPeerConn(Conn)
	GetServerConn() Conn
	SetServerConn(Conn)
	Connect()
	Stop()
	Start() error
	RegisteredCallback(Client)
	ConnectingCallback(Client)
	ConnectedCallback(Client)
	MessageCallback(Client, string)
	OnRegistered(func(Client))
	OnConnecting(func(Client))
	OnConnected(func(Client))
	OnMessage(func(Client, string))
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

type History struct {
	w *bufio.Writer
	m *sync.Mutex
}

func (h *History) Add(text string) {
	h.m.Lock()
	defer h.m.Unlock()
	fmt.Fprintln(h.w, text)
	h.w.Flush()
}

func NewHistory(f *os.File) *History {
	return &History{
		w: bufio.NewWriter(f),
		m: &sync.Mutex{},
	}
}

type Message struct {
	Type    string      `json:"type"`
	PeerID  string      `json:"peerID, omitempty"`
	Error   string      `json:"error, omitempty"`
	Content interface{} `json:"data, omitempty"`
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
	PublicKey  string       `json:"publickKey,omitempty"`
	PrivateKey string       `json:"-"`
	Addr       *net.UDPAddr `json:"-"`
}

func (p *Peer) GetPublicKey() ([32]byte, error) {
	var key [32]byte
	bs, err := base64.StdEncoding.DecodeString(p.PublicKey)
	if err != nil {
		return key, err
	}
	copy(key[:], bs)
	return key, nil
}

func (p *Peer) SetPublicKey(key [32]byte) {
	p.PublicKey = base64.StdEncoding
}

type Peers map[string]*Peer
