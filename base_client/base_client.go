package base_client

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"p2p-network-test/shared"
	"strconv"
	"sync"
)

type Client struct {
	s                  shared.Server
	self               *shared.Peer
	log                *log.Logger
	logFile            *os.File
	peer               *shared.Peer
	sConn              shared.Conn
	pConn              shared.Conn
	keySent            bool
	keyReceived        bool
	mKeySent           *sync.RWMutex
	mKeyReceived       *sync.RWMutex
	mPConn             *sync.Mutex
	resetCallback      func(shared.Client)
	registeredCallback func(shared.Client)
	connectingCallback func(shared.Client)
	connectedCallback  func(shared.Client)
	messageCallback    func(shared.Client, string)
}

func (c *Client) WasKeySent() bool {
	c.mKeySent.RLock()
	defer c.mKeySent.RUnlock()
	return c.keySent
}

func (c *Client) WasKeyReceived() bool {
	c.mKeyReceived.RLock()
	defer c.mKeyReceived.RUnlock()
	return c.keyReceived
}

func (c *Client) SetKeySent(b bool) {
	c.mKeySent.Lock()
	defer c.mKeySent.Unlock()
	c.keySent = b
}

func (c *Client) SetKeyReceived(b bool) {
	c.mKeyReceived.Lock()
	defer c.mKeyReceived.Unlock()
	c.mKeyReceived = b
}

func (c *Client) GetServer() shared.Server {
	return c.s
}

func (c *Client) GetLog() *log.Logger {
	return c.log
}

func (c *Client) GetSelf() *shared.Peer {
	return c.self
}

func (c *Client) GetPeer() *shared.Peer {
	return c.peer
}

func (c *Client) SetPeer(p *shared.Peer) {
	c.peer = p
}

func (c *Client) GetPeerConn() shared.Conn {
	c.mPConn.Lock()
	defer c.mPConn.Unlock()
	return c.pConn
}

func (c *Client) SetPeerConn(c shared.Conn) {
	c.mPConn.Lock()
	defer c.mPConn.Unlock()
	c.pConn = c
}

func (c *Client) GetServerConn() shared.Conn {
	return c.sConn
}

func (c *Client) SetServerConn(c shared.Conn) {
	c.sConn = c
}

func (c *Client) Stop() {
	c.s.Stop()
}

func (c *Client) RegisteredCallback(client shared.Client) {
	c.registeredCallback(client)
}

func (c *Client) ConnectingCallback(client shared.Client) {
	c.connectingCallback(client)
}

func (c *Client) ConnectedCallback(client shared.Client) {
	c.connectedCallback(client)
}

func (c *Client) MessageCallback(client shared.Client) {
	c.messageCallback(client, text)
}

func (c *Client) OnReset(f func(shared.Client)) {
	c.resetCallback = f
}

func (c *Client) OnRegistered(f func(shared.Client)) {
	c.registeredCallback = f
}

func (c *Client) OnConnecting(f func(shared.Client)) {
	c.connectingCallback = f
}

func (c *Client) OnConnected(f func(shared.Client)) {
	c.connectedCallback = f
}

func (c *Client) OnMessage(f func(shared.Client, string)) {
	c.messageCallback = f
}

func New(username string, s shared.Server) (*Client, error) {
	self := &shared.Peer{Username: username}

	self.ID = strconv.Itoa(rand.Intn(65535))

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	lf, err := os.Create(fmt.Sprintf("%s/log-%s.txt", wd, self.Username))
	if err != nil {
		return nil, err
	}

	l := log.New(lf, "", log.LstdFlags|log.Lshortfile)
	l.Printf("Logging initialized")

	p := &shared.Peer{}

	return &Client{
		s:                  s,
		self:               self,
		peer:               p,
		log:                l,
		logFile:            lf,
		mKeyReceived:       &sync.RWMutex{},
		mKeySent:           &sync.RWMutex{},
		mPConn:             &sync.Mutex{},
		resetCallback:      func(shared.Client) {},
		registeredCallback: func(shared.Client) {},
		connectingCallback: func(shared.Client) {},
		connectedCallback:  func(shared.Client) {},
		messageCallback:    func(shared.Client, string) {},
	}, nil
}
