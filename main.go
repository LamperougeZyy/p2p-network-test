package main

import (
	"fmt"
	"log"
	"net"

	"p2p-network-test/shared"
	"p2p-network-test/udp_server"
)

//var (
//	pubKey [32]byte
//	priKey [32]byte
//)

func route(peers shared.Peers, cs shared.Conns, c shared.Conn, m *shared.Message) (*shared.Message, error) {
	//根据消息类型选择处理函数
	switch m.Type {
	case "greeting":
		return greetingHandler(c, m)
	case "register":
		return registerHandler(peers, c, m)
	case "establish":
		return establishHandler(peers, cs, m)
	default:
		return notFoundHandler(m)
	}
}

func createMessageCallBack(peers shared.Peers) func(cs shared.Conns, c shared.Conn, m *shared.Message) {
	return func(cs shared.Conns, c shared.Conn, m *shared.Message) {
		//返回客户端的ip地址，所采用的协议，以及消息类型
		log.Printf("Request from client at %s over %s with type %s", c.GetAddr(), c.Protocol(), m.Type)

		//处理消息，找到相应的句柄
		res, err := route(peers, cs, c, m)

		//出错时返回错误消息
		if err != nil {
			c.Send(&shared.Message{
				Type:  m.Type,
				Error: err.Error(),
			})
			return
		}

		//正确返回
		err = c.Send(res)
		if err != nil {
			log.Print(err)
		}
	}
}

func main() {
	fmt.Println("Start my own p2p network by UDP hole punching")

	//将地址作为UDP地址进行解析并返回
	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:9000")
	if err != nil {
		log.Fatal(err)
	}

	//服务端实例化
	udpS, err := udp_server.New(udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	udpPeers := make(shared.Peers)
	udpS.OnMessage(createMessageCallBack(udpPeers)) //设置服务端消息回调函数
	udpS.Listen()
}
