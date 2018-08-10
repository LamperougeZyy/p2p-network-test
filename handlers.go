package main

import (
	//"encoding/base64"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"p2p-network-test/shared"
	"strconv"
	"strings"
)

func greetingHandler(conn shared.Conn, m *shared.Message) (*shared.Message, error) {
	//检测消息中是否带有公钥
	//str, ok := m.Content.(string)
	//if !ok {
	//	return nil, errors.New("greeting request must contain client's public key")
	//}

	//
	//bs, err := base64.StdEncoding.DecodeString(str)
	//if err != nil {
	//	return nil, err
	//}

	return &shared.Message{
		Type:    "greeting",
		Content: "Hello",
	}, nil
}

func registerHandler(peers shared.Peers, c shared.Conn, m *shared.Message) (*shared.Message, error) {
	var registration shared.Registration
	err := mapstructure.Decode(m.Content, &registration) //取出m.Content中包含registration字段的内容
	if err != nil {
		return nil, err
	}

	//将该节点添加到已注册列表中
	endpoint := strings.Split(c.GetAddr().String(), ":")
	if len(endpoint) != 2 {
		return nil, errors.New("address is not valid")
	}

	port, err := strconv.Atoi(endpoint[1])
	if err != nil {
		return nil, err
	}

	peers[m.PeerID] = &shared.Peer{
		ID:       m.PeerID,
		Username: registration.Username,
		Endpoint: shared.Endpoint{
			IP:   endpoint[0],
			Port: port,
		},
	}
	log.Printf("Registered peer: %s at addr %s", m.PeerID, c.GetAddr().String())

	return &shared.Message{
		Type:    "register",
		Encrypt: true,
	}, nil
}

func establishHandler(peers shared.Peers, conns shared.Conns, m *shared.Message) (*shared.Message, error) {
	// 验证被请求节点是否已经注册
	rp, ok := peers[m.PeerID]
	if !ok {
		return nil, errors.New("client is not registered with this server")
	}

	// 验证请求的消息是否有效
	id, ok := m.Content.(string)
	if !ok {
		return nil, errors.New("request content is malformed")
	}

	// make sure the other peer has registered with the server
	op, ok := peers[id]
	if !ok {
		return nil, fmt.Errorf("The peer: %s has not registered with the server.", id)
	}

	// get conn for other peer
	conn, ok := conns[op.Endpoint.String()]
	if !ok {
		return nil, fmt.Errorf("Could not resolve the peer: %s's conn", id)
	}

	// send requesting peer's endpoint to other peer
	conn.Send(&shared.Message{
		Type:    "establish",
		Content: rp,
		Encrypt: true,
	})

	// send requesting peer other peer's endpoint
	return &shared.Message{
		Type:    "establish",
		Content: op,
		Encrypt: true,
	}, nil
}

func notFoundHandler(m *shared.Message) (*shared.Message, error) {
	return nil, fmt.Errorf("Request type %s undefined", m.Type)
}
